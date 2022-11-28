package repository

import (
	"fmt"
	tableConstants "main-server/pkg/constant/table"
	userModel "main-server/pkg/model/user"

	"github.com/jmoiron/sqlx"
)

type AuthTypePostgres struct {
	db *sqlx.DB
}

/*
* Функция создания экземпляра сервиса
 */
func NewAuthTypePostgres(db *sqlx.DB) *AuthTypePostgres {
	return &AuthTypePostgres{db: db}
}

/*
* Функция получения данных о роли
 */
func (r *AuthTypePostgres) GetAuthType(column, value interface{}) (userModel.AuthTypeModel, error) {
	var data userModel.AuthTypeModel
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", tableConstants.AUTH_TYPES_TABLE, column.(string))

	var err error

	switch value.(type) {
	case int:
		err = r.db.Get(&data, query, value.(int))
		break
	case string:
		err = r.db.Get(&data, query, value.(string))
		break
	}

	return data, err
}
