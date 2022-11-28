package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type WrapperPostgres struct {
	db *sqlx.DB
}

func NewWrapperPostgres(db *sqlx.DB) *WrapperPostgres {
	return &WrapperPostgres{
		db: db,
	}
}

func (r *WrapperPostgres) GetOne(table string, column, value interface{}) (interface{}, error) {
	var model interface{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1 LIMIT 1", table, column.(string))

	var err error

	switch value.(type) {
	case int:
		err = r.db.Get(&model, query, value.(int))
		break
	case string:
		err = r.db.Get(&model, query, value.(string))
		break
	}

	return model, err
}
