package repository

import (
	//middlewareConstant "main-server/pkg/constant/middleware"
	"encoding/json"
	"fmt"
	actionConstant "main-server/pkg/constant/action"
	middlewareConstant "main-server/pkg/constant/middleware"
	objectConstant "main-server/pkg/constant/object"
	roleConstant "main-server/pkg/constant/role"
	tableConstant "main-server/pkg/constant/table"
	adminModel "main-server/pkg/model/admin"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	"main-server/pkg/model/worker"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type AdminPostgres struct {
	db       *sqlx.DB
	enforcer *casbin.Enforcer
	domain   *DomainPostgres
	role     *RolePostgres
}

/* Function for create new struct of AdminPostgres */
func NewAdminPostgres(
	db *sqlx.DB,
	enforcer *casbin.Enforcer,
	domain *DomainPostgres,
	role *RolePostgres,
) *AdminPostgres {
	return &AdminPostgres{
		db:       db,
		enforcer: enforcer,
		domain:   domain,
		role:     role,
	}
}

/* Метод для получения информации обо всех пользователях */
func (r *AdminPostgres) GetAllUsers(c *gin.Context) (adminModel.UsersResponseModel, error) {
	// Получение email-адресов всех пользователей
	query := fmt.Sprintf(`SELECT email FROM %s`, tableConstant.USERS_TABLE)
	var users []adminModel.UserResponseModel

	if err := r.db.Select(&users, query); err != nil {
		return adminModel.UsersResponseModel{}, err
	}

	return adminModel.UsersResponseModel{
		Users: &users,
	}, nil
}

/* Функция создания новой компании (доступно админу) */
func (r *AdminPostgres) CreateCompany(c *gin.Context, data adminModel.CompanyModel) (adminModel.CompanyModel, error) {
	usersId, _ := c.Get(middlewareConstant.USER_CTX)
	domainsId, _ := c.Get(middlewareConstant.DOMAINS_ID)

	// Обёртка всех операций в одну транзакцию в рамках текущей функции
	tx, err := r.db.Begin()
	if err != nil {
		return adminModel.CompanyModel{}, err
	}

	query := fmt.Sprintf("INSERT INTO %s (uuid, data, created_at, updated_at, users_id) values ($1, $2, $3, $4, $5) RETURNING id", tableConstant.COMPANIES_TABLE)

	// Преобразование данных структуры в данные JSON
	dataJson, err := json.Marshal(data)
	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	var companyId int
	currentDate := time.Now()
	companyUuid := uuid.NewV4()

	row := tx.QueryRow(query, companyUuid, dataJson, currentDate, currentDate, usersId)
	if err := row.Scan(&companyId); err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	// Добавление информации о новом объекте (объект в данном случае - это компания)
	var typesObjects rbacModel.TypesObjectsModel

	query = fmt.Sprintf("SELECT * FROM %s WHERE value=$1", tableConstant.TYPES_OBJECTS_TABLE)

	err = r.db.Get(&typesObjects, query, objectConstant.TYPE_COMPANY)
	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	query = fmt.Sprintf("INSERT INTO %s (value, types_objects_id) values ($1, $2)", tableConstant.OBJECTS_TABLE)

	_, err = tx.Exec(query, companyUuid, typesObjects.Id)
	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	// Добавление информации о пользователе в текущей компании
	var user userModel.UserModel
	query = fmt.Sprintf("SELECT * FROM %s tl WHERE tl.email=$1", tableConstant.USERS_TABLE)
	err = r.db.Get(&user, query, data.EmailAdmin)
	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	query = fmt.Sprintf(`
	INSERT INTO %s (uuid, data, created_at, updated_at, users_id, companies_id) 
	values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		tableConstant.WORKERS_TABLE,
	)

	roleAdmin, err := r.role.GetRole("value", roleConstant.ROLE_BUILDER_ADMIN)
	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	roleAdminJson, err := json.Marshal(worker.WorkerModel{
		Role: roleAdmin.Uuid,
	})

	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	var workerId int
	row = tx.QueryRow(query, uuid.NewV4().String(), roleAdminJson, currentDate, currentDate, user.Id, companyId)
	if err := row.Scan(&workerId); err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	var userId string = strconv.Itoa(usersId.(int))
	var domainId string = strconv.Itoa(domainsId.(int))
	var roleAdminId string = strconv.Itoa(roleAdmin.Id)

	// Обновление политик управления доступа для текущего пользователя
	_, err = r.enforcer.AddPolicies([][]string{
		{userId, domainId, companyUuid.String(), actionConstant.DELETE},
		{userId, domainId, companyUuid.String(), actionConstant.MODIFY},
		{userId, domainId, companyUuid.String(), actionConstant.READ},
		{userId, domainId, companyUuid.String(), actionConstant.ADMINISTRATION},
	})

	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	_, err = r.enforcer.AddRoleForUserInDomain(userId, roleAdminId, domainId)

	if err != nil {
		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	// Get id for user admin
	query = fmt.Sprintf("SELECT id FROM %s WHERE email=$1 LIMIT 1", tableConstant.USERS_TABLE)
	var userAdminId int

	row = r.db.QueryRow(query, data.EmailAdmin)
	if err := row.Scan(&userAdminId); err != nil {
		r.enforcer.RemovePolicies([][]string{
			{userId, domainId, companyUuid.String(), actionConstant.DELETE},
			{userId, domainId, companyUuid.String(), actionConstant.MODIFY},
			{userId, domainId, companyUuid.String(), actionConstant.READ},
			{userId, domainId, companyUuid.String(), actionConstant.ADMINISTRATION},
		})

		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	userId = strconv.Itoa(userAdminId)

	// For admin company
	_, err = r.enforcer.AddPolicies([][]string{
		{userId, domainId, companyUuid.String(), actionConstant.MODIFY},
		{userId, domainId, companyUuid.String(), actionConstant.READ},
		{userId, domainId, companyUuid.String(), actionConstant.ADMINISTRATION},
	})

	// Save results all operation into a tables
	if err := tx.Commit(); err != nil {
		r.enforcer.RemovePolicies([][]string{
			{userId, domainId, companyUuid.String(), actionConstant.MODIFY},
			{userId, domainId, companyUuid.String(), actionConstant.READ},
			{userId, domainId, companyUuid.String(), actionConstant.ADMINISTRATION},
		})

		tx.Rollback()
		return adminModel.CompanyModel{}, err
	}

	return data, nil
}
