package project

import (
	"time"
)

/* Model data for response users model */
type ProjectModel struct {
	Logo        *string          `json:"logo"`
	Uuid        string           `json:"uuid" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Managers    []UserEmailModel `json:"managers" binding:"required"`
}

/* Model data for request update project in company */
type ProjectUpdateModel struct {
	Uuid        string           `json:"uuid" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Managers    []UserEmailModel `json:"managers" binding:"required"`
}

type ProjectImageModel struct {
	Filepath string `json:"filepath" binding:"required" db:"filepath"`
	Uuid     string `json:"uuid" binding:"required"`
}

type ProjectDataModel struct {
	Logo        *string          `json:"logo" db:"logo"`
	Title       string           `json:"title" binding:"required" db:"title"`
	Description string           `json:"description" binding:"required" db:"description"`
	Managers    []UserEmailModel `json:"managers" binding:"required" db:"managers"`
}

type ProjectDbModel struct {
	Uuid      string    `json:"uuid" db:"uuid"`
	Data      string    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ProjectDbDataEx struct {
	Uuid      string           `json:"uuid" db:"uuid"`
	Data      ProjectDataModel `json:"data" db:"data"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

type ProjectUuidModel struct {
	Uuid string `json:"uuid" binding:"required"`
}

type ProjectCountModel struct {
	Uuid  string `json:"uuid" binding:"required"`
	Limit int    `json:"limit" binding:"required"`
	Count int    `json:"count"`
}

type ProjectAnyCountModel struct {
	Projects []ProjectDbDataEx `json:"projects" binding:"required"`
	Count    int               `json:"count" binding:"required"`
}
