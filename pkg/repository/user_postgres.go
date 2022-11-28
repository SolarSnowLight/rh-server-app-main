package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	actionConstant "main-server/pkg/constant/action"
	middlewareConstant "main-server/pkg/constant/middleware"
	objectConstant "main-server/pkg/constant/object"
	tableConstant "main-server/pkg/constant/table"
	companyModel "main-server/pkg/model/company"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserPostgres struct {
	db       *sqlx.DB
	enforcer *casbin.Enforcer
	domain   *DomainPostgres
	role     *RolePostgres
}

/*
* Функция создания экземпляра сервиса
 */
func NewUserPostgres(
	db *sqlx.DB, enforcer *casbin.Enforcer,
	domain *DomainPostgres, role *RolePostgres,
) *UserPostgres {
	return &UserPostgres{
		db:       db,
		enforcer: enforcer,
		domain:   domain,
		role:     role,
	}
}

func (r *UserPostgres) GetUser(column, value interface{}) (userModel.UserModel, error) {
	var user userModel.UserModel
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", tableConstant.USERS_TABLE, column.(string))

	var err error

	switch value.(type) {
	case int:
		err = r.db.Get(&user, query, value.(int))
		break
	case string:
		err = r.db.Get(&user, query, value.(string))
		break
	}

	return user, err
}

func (r *UserPostgres) GetProfile(c *gin.Context) (userModel.UserProfileModel, error) {
	usersId, _ := c.Get(middlewareConstant.USER_CTX)

	var profile userModel.UserProfileModel
	var email userModel.UserEmailModel

	query := fmt.Sprintf("SELECT data FROM %s tl WHERE tl.users_id = $1 LIMIT 1",
		tableConstant.USERS_DATA_TABLE,
	)

	err := r.db.Get(&profile, query, usersId)
	if err != nil {
		return userModel.UserProfileModel{}, err
	}

	query = fmt.Sprintf("SELECT email FROM %s tl WHERE tl.id = $1 LIMIT 1", tableConstant.USERS_TABLE)

	err = r.db.Get(&email, query, usersId)
	if err != nil {
		return userModel.UserProfileModel{}, err
	}

	return userModel.UserProfileModel{
		Email: email.Email,
		Data:  profile.Data,
	}, nil
}

func (r *UserPostgres) UpdateProfile(c *gin.Context, data userModel.UserProfileUpdateDataModel) (userModel.UserJSONBModel, error) {
	usersId, _ := c.Get(middlewareConstant.USER_CTX)

	userJsonb, err := json.Marshal(data)
	if err != nil {
		return userModel.UserJSONBModel{}, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return userModel.UserJSONBModel{}, err
	}

	query := fmt.Sprintf("UPDATE %s tl SET data=$1 WHERE tl.users_id = $2", tableConstant.USERS_DATA_TABLE)

	// Update data about user profile
	_, err = tx.Exec(query, userJsonb, usersId)
	if err != nil {
		tx.Rollback()
		return userModel.UserJSONBModel{}, err
	}

	query = fmt.Sprintf("SELECT data FROM %s tl WHERE users_id=$1 LIMIT 1", tableConstant.USERS_DATA_TABLE)
	var userData []userModel.UserDataModel

	err = r.db.Select(&userData, query, usersId)

	if err != nil {
		tx.Rollback()
		return userModel.UserJSONBModel{}, err
	}

	if len(userData) <= 0 {
		tx.Rollback()
		return userModel.UserJSONBModel{}, errors.New("Данных у пользователя нет")
	}

	var dataFromJson userModel.UserJSONBModel
	err = json.Unmarshal([]byte(userData[0].Data), &dataFromJson)

	if err != nil {
		tx.Rollback()
		return userModel.UserJSONBModel{}, err
	}

	// Change password
	if data.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), viper.GetInt("crypt.cost"))
		if err != nil {
			tx.Rollback()
			return userModel.UserJSONBModel{}, err
		}

		query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE id=$2", tableConstant.USERS_TABLE)
		_, err = r.db.Exec(query, string(hashedPassword), usersId)

		if err != nil {
			tx.Rollback()
			return userModel.UserJSONBModel{}, err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return userModel.UserJSONBModel{}, err
	}

	return dataFromJson, nil
}

/* Get information about company for current user */
func (r *UserPostgres) GetUserCompany(userId, domainId int) (companyModel.CompanyDbModelEx, error) {
	query := fmt.Sprintf(`
			SELECT DISTINCT tl.value, trule.v3
			FROM %s tl
			JOIN %s tr ON tr.id = tl.types_objects_id
			JOIN %s trule ON trule.v2 = tl.value 
			WHERE tr.value = $1 
					AND trule.v0 = $2 
					AND trule.ptype = 'p'
					AND (trule.v3 = $3 OR trule.v3 = $4)
	`, tableConstant.OBJECTS_TABLE, tableConstant.TYPES_OBJECTS_TABLE, tableConstant.RULES_TABLE)

	var companies []companyModel.CompanyRuleModelEx
	err := r.db.Select(&companies, query,
		objectConstant.TYPE_COMPANY, userId,
		actionConstant.ADMINISTRATION, actionConstant.MANAGEMENT,
	)

	if err != nil {
		return companyModel.CompanyDbModelEx{}, err
	}

	if len(companies) <= 0 {
		return companyModel.CompanyDbModelEx{}, nil
	}

	company := companies[len(companies)-1]
	policies := r.enforcer.GetFilteredPolicy(0, strconv.Itoa(userId), strconv.Itoa(domainId), company.Value)

	rules := lo.Map(policies, func(x []string, _ int) string {
		return x[len(x)-1]
	})

	var companyInfo companyModel.CompanyDbModel
	query = fmt.Sprintf(`SELECT uuid, data, created_at FROM %s tl WHERE tl.uuid = $1 LIMIT 1`, tableConstant.COMPANIES_TABLE)

	err = r.db.Get(&companyInfo, query, company.Value)
	if err != nil {
		return companyModel.CompanyDbModelEx{}, err
	}

	var companyDataEx companyModel.CompanyModel
	err = json.Unmarshal([]byte(companyInfo.Data), &companyDataEx)
	if err != nil {
		return companyModel.CompanyDbModelEx{}, err
	}

	return companyModel.CompanyDbModelEx{
		Uuid:      companyInfo.Uuid,
		Data:      companyDataEx,
		CreatedAt: companyInfo.CreatedAt,
		Rules:     rules,
	}, nil
}

/* Method for to check acces every user. On result this method make decision for make navbar or other component */
func (r *UserPostgres) AccessCheck(userId, domainId int, value rbacModel.RoleValueModel) (bool, error) {
	role, err := r.role.GetRole("value", value.Value)

	if err != nil {
		return false, err
	}

	result, err := r.enforcer.HasRoleForUser(strconv.Itoa(userId), strconv.Itoa(role.Id), strconv.Itoa(domainId))

	if err != nil {
		return false, err
	}

	return result, nil
}
