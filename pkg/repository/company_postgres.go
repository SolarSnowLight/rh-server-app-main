package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	actionConstant "main-server/pkg/constant/action"
	roleConstant "main-server/pkg/constant/role"
	tableConstant "main-server/pkg/constant/table"
	companyModel "main-server/pkg/model/company"
	projectModel "main-server/pkg/model/project"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	utilSlice "main-server/pkg/util"
	"sort"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
)

type CompanyPostgres struct {
	db       *sqlx.DB
	enforcer *casbin.Enforcer
	role     *RolePostgres
	user     *UserPostgres
	wrapper  *WrapperPostgres
}

/* Function for create new struct of CompanyPostgres */
func NewCompanyPostgres(
	db *sqlx.DB,
	enforcer *casbin.Enforcer,
	role *RolePostgres,
	user *UserPostgres,
	wrapper *WrapperPostgres,
) *CompanyPostgres {
	return &CompanyPostgres{
		db:       db,
		enforcer: enforcer,
		role:     role,
		user:     user,
		wrapper:  wrapper,
	}
}

/* Получение списка менеджеров компании */
func (r *CompanyPostgres) GetManagers(userId, domainId int, data companyModel.ManagerCountModel) (companyModel.ManagerAnyCountModel, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE uuid=$1", tableConstant.COMPANIES_TABLE)
	var companyId int

	row := r.db.QueryRow(query, data.Uuid)
	if err := row.Scan(&companyId); err != nil {
		return companyModel.ManagerAnyCountModel{}, err
	}

	var managers []companyModel.ManagerDbDataEx
	sum := (data.Count + data.Limit)

	query = fmt.Sprintf(`
	SELECT DISTINCT tl.uuid, tl.email, tr.data, tw.created_at FROM %s tl
	JOIN %s tr ON tr.users_id = tl.id
	JOIN %s tw ON tw.users_id = tl.id
	JOIN %s twp ON twp.workers_id = tw.id
	WHERE tw.data = $1 AND tw.companies_id = $2
	`,
		tableConstant.USERS_TABLE, tableConstant.USERS_DATA_TABLE,
		tableConstant.WORKERS_TABLE, tableConstant.WORKERS_PROJECTS_TABLE,
	)

	role, err := r.role.GetRole("value", roleConstant.ROLE_BUILDER_MANAGER)
	if err != nil {
		return companyModel.ManagerAnyCountModel{}, err
	}

	roleData := rbacModel.RoleDataModel{Role: role.Uuid}
	roleDataStr, err := json.Marshal(roleData)
	if err != nil {
		return companyModel.ManagerAnyCountModel{}, err
	}

	err = r.db.Select(&managers, query, roleDataStr, companyId)
	if err != nil {
		return companyModel.ManagerAnyCountModel{}, err
	}

	var managersEx []companyModel.ManagerDataEx

	for _, element := range managers {
		var managerData companyModel.ManagerUserData
		err = json.Unmarshal([]byte(element.Data), &managerData)

		if err != nil {
			return companyModel.ManagerAnyCountModel{}, err
		}

		managersEx = append(managersEx, companyModel.ManagerDataEx{
			Uuid:      element.Uuid,
			Email:     element.Email,
			Data:      managerData,
			CreatedAt: element.CreatedAt,
		})
	}

	// Remove duplicates
	managersEx = utilSlice.RemoveDuplicate(managersEx)

	sort.SliceStable(managersEx, func(i, j int) bool {
		return managersEx[i].CreatedAt.After(managersEx[j].CreatedAt)
	})

	if data.Count >= len(managersEx) {
		return companyModel.ManagerAnyCountModel{}, nil
	}

	if sum >= len(managersEx) {
		sum -= (sum - len(managersEx))
	}

	return companyModel.ManagerAnyCountModel{
		Managers: managersEx[data.Count:sum],
		Count:    (sum - data.Count),
	}, nil
}

/* Method for update image for current company */
func (r *CompanyPostgres) CompanyUpdateImage(user userModel.UserIdentityModel, data companyModel.CompanyImageModel) (companyModel.CompanyImageModel, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return companyModel.CompanyImageModel{}, err
	}

	userIdStr := strconv.Itoa(user.UserId)
	domainIdStr := strconv.Itoa(user.DomainId)

	// Check access for user
	access, err := r.enforcer.Enforce(userIdStr, domainIdStr, data.Uuid, actionConstant.MODIFY)
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyImageModel{}, err
	}

	if !access {
		tx.Rollback()
		return companyModel.CompanyImageModel{}, errors.New("Ошибка! Нет доступа!")
	}

	// Update logo for project
	query := fmt.Sprintf(`UPDATE %s tl SET data = jsonb_set(data, '{logo}', to_jsonb($1::text), true) WHERE tl.uuid = $2`, tableConstant.COMPANIES_TABLE)

	_, err = r.db.Exec(query, data.Filepath, data.Uuid)
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyImageModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyImageModel{}, err
	}

	return data, nil
}

/* Method for update info for current company */
func (r *CompanyPostgres) CompanyUpdate(user userModel.UserIdentityModel, data companyModel.CompanyUpdateModel) (companyModel.CompanyUpdateModel, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return companyModel.CompanyUpdateModel{}, err
	}

	userIdStr := strconv.Itoa(user.UserId)
	domainIdStr := strconv.Itoa(user.DomainId)

	// Check access for user
	access, err := r.enforcer.Enforce(userIdStr, domainIdStr, data.Uuid, actionConstant.MODIFY)
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyUpdateModel{}, err
	}

	if !access {
		tx.Rollback()
		return companyModel.CompanyUpdateModel{}, errors.New("Ошибка! Нет доступа!")
	}

	var companyInfo []companyModel.CompanyDbModel
	query := fmt.Sprintf("SELECT uuid, data, created_at FROM %s WHERE uuid=$1 LIMIT 1", tableConstant.COMPANIES_TABLE)

	err = r.db.Select(&companyInfo, query, data.Uuid)
	if err != nil {
		return companyModel.CompanyUpdateModel{}, err
	}

	var companyData companyModel.CompanyModel
	err = json.Unmarshal([]byte(companyInfo[0].Data), &companyData)

	companyData.Description = data.Description
	companyData.EmailCompany = data.EmailCompany
	companyData.Link = data.Link
	companyData.Phone = data.Phone
	companyData.Title = data.Title

	companyDataJson, err := json.Marshal(companyData)
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyUpdateModel{}, err
	}

	query = fmt.Sprintf("UPDATE %s tl SET data=$1 WHERE tl.uuid=$2", tableConstant.COMPANIES_TABLE)

	_, err = tx.Exec(query, companyDataJson, data.Uuid)
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyUpdateModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return companyModel.CompanyUpdateModel{}, err
	}

	return data, nil
}

type SQLResultGetManager struct {
	Id int `json:"id" binding:"required" db:"id"`
}

type SQLResultListProject struct {
	Uuid string `json:"uuid" binding:"required" db:"uuid"`
	Data string `json:"data" binding:"required" db:"data"`
}

/* Метод структуры для получения информации о менеджере */
func (r *CompanyPostgres) GetManager(user userModel.UserIdentityModel, data companyModel.ManagerUuidModel) (companyModel.ManagerCompanyModel, error) {
	// Получение информации о пользователе по UUID
	manager, err := r.user.GetUser("uuid", data.ManagerUuid)
	if err != nil {
		return companyModel.ManagerCompanyModel{}, err
	}

	// Информация о менеджере
	var sqlResultManager []SQLResultGetManager

	query := fmt.Sprintf(
		`SELECT DISTINCT w.id FROM %s u
		INNER JOIN %s ud ON ud.users_id = u.id
		INNER JOIN %s w ON w.users_id = u.id
		WHERE u.id = $1 LIMIT 1`,
		tableConstant.USERS_TABLE,
		tableConstant.USERS_DATA_TABLE,
		tableConstant.WORKERS_TABLE,
	)

	if err := r.db.Select(&sqlResultManager, query, manager.Id); err != nil {
		return companyModel.ManagerCompanyModel{}, err
	}

	// Проверка полученного массива данных
	if len(sqlResultManager) <= 0 {
		return companyModel.ManagerCompanyModel{}, nil
	}

	// Информация о всех проектах, в которых участвует менеджер
	var sqlResultProject []SQLResultListProject

	query = fmt.Sprintf(
		`SELECT p.uuid, p.data FROM %s wp
		INNER JOIN %s p ON p.id = wp.projects_id
		WHERE wp.workers_id = $1;`,
		tableConstant.WORKERS_PROJECTS_TABLE,
		tableConstant.PROJECTS_TABLE,
	)

	if err := r.db.Select(&sqlResultProject, query, sqlResultManager[len(sqlResultManager)-1].Id); err != nil {
		return companyModel.ManagerCompanyModel{}, err
	}

	var projects []companyModel.ManagerProjectInfoModel

	for _, element := range sqlResultProject {
		var data projectModel.ProjectDataModel
		err := json.Unmarshal([]byte(element.Data), &data)
		if err != nil {
			return companyModel.ManagerCompanyModel{}, err
		}

		projects = append(projects, companyModel.ManagerProjectInfoModel{
			Uuid:        element.Uuid,
			Logo:        data.Logo,
			Title:       data.Title,
			Description: data.Description,
		})
	}

	return companyModel.ManagerCompanyModel{
		Projects: projects,
	}, nil
}
