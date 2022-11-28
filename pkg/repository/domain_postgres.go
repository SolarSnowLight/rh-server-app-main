package repository

import (
	"fmt"
	tableConstants "main-server/pkg/constant/table"
	rbacModel "main-server/pkg/model/rbac"

	"github.com/jmoiron/sqlx"
)

type DomainPostgres struct {
	db *sqlx.DB
}

/*
* Функция создания экземпляра сервиса
 */
func NewDomainPostgres(db *sqlx.DB) *DomainPostgres {
	return &DomainPostgres{db: db}
}

/* Get information about domain */
func (r *DomainPostgres) GetDomain(column, value interface{}) (rbacModel.DomainModel, error) {
	var domain rbacModel.DomainModel
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", tableConstants.DOMAINS_TABLE, column.(string))

	var err error

	switch value.(type) {
	case int:
		err = r.db.Get(&domain, query, value.(int))
		break
	case string:
		err = r.db.Get(&domain, query, value.(string))
		break
	}

	return domain, err
}
