package article

import "time"

/* Model data for request create article */
type ArticleCreateRequestModel struct {
	Title    string                  `json:"title" binding:"required"`
	Text     string                  `json:"text" binding:"required"`
	Filename *string                 `json:"filename" binding:"required"`
	Filepath *string                 `json:"filepath" binding:"required"`
	Tags     string                  `json:"tags" binding:"required"`
	Files    *[]ArticlesFilesDBModel `json:"files" binding:"required"`
}

/* Model data for request update article */
type ArticleUpdateRequestModel struct {
	Uuid        string                  `json:"uuid" binding:"required" db:"uuid"`
	Title       string                  `json:"title" binding:"required"`
	Text        string                  `json:"text" binding:"required"`
	Filename    *string                 `json:"filename" binding:"required"`
	Filepath    *string                 `json:"filepath" binding:"required"`
	Tags        string                  `json:"tags" binding:"required"`
	Files       *[]ArticlesFilesDBModel `json:"files" binding:"required"`
	FilesDelete *[]int                  `json:"files_delete" binding:"required"`
}

type FileArticleExModel struct {
	Filename string
	Filepath string
	Index    int
	Id       int
}

type ArticleSuccessModel struct {
	Success bool `json:"success"`
}

type ArticleModel struct {
	Uuid      string                 `json:"uuid" binding:"required"`
	Filepath  string                 `json:"filepath" binding:"required"`
	Title     string                 `json:"title" binding:"required"`
	Text      string                 `json:"text" binding:"required"`
	Tags      string                 `json:"tags" binding:"required"`
	Files     []ArticlesFilesDBModel `json:"files" binding:"required"`
	CreatedAt time.Time              `json:"created_at" binding:"required"`
	UpdatedAt time.Time              `json:"updated_at" binding:"required"`
}

type ArticlesModel struct {
	Articles []ArticleModel `json:"articles" binding:"required"`
}

type ArticleUuidModel struct {
	Uuid string `json:"uuid" binding:"required"`
}

type ArticleDBModel struct {
	Id        int       `json:"id" binding:"required" db:"id"`
	Uuid      string    `json:"uuid" binding:"required" db:"uuid"`
	UsersId   int       `json:"users_id" binding:"required" db:"users_id"`
	Filepath  string    `json:"filepath" binding:"required" db:"filepath"`
	Filename  string    `json:"filename" binding:"required" db:"filename"`
	Title     string    `json:"title" binding:"required" db:"title"`
	Text      string    `json:"text" binding:"required" db:"text"`
	Tags      string    `json:"tags" binding:"required" db:"tags"`
	CreatedAt time.Time `json:"created_at" binding:"required" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" binding:"required" db:"updated_at"`
}

type ArticlesFilesDBModel struct {
	FilesId  *int   `json:"files_id" db:"files_id"`
	Index    int    `json:"index" binding:"required" db:"index"`
	Filename string `json:"filename" binding:"required" db:"filename"`
	Filepath string `json:"filepath" binding:"required" db:"filepath"`
}

type ArticlesFilesModel struct {
	FilesId    int `json:"files_id" db:"files_id"`
	Index      int `json:"index" binding:"required" db:"index"`
	ArticlesId int `json:"articles_id" binding:"required" db:"articles_id"`
	Id         int `json:"id" binding:"required" db:"id"`
}

/* Common data model for article */
type ArticleDataModel struct {
	Uuid      string                 `json:"uuid" binding:"required"`
	Title     string                 `json:"title" binding:"required"`
	Text      string                 `json:"text" binding:"required"`
	Filename  string                 `json:"filename" binding:"required"`
	Filepath  string                 `json:"filepath" binding:"required"`
	Tags      string                 `json:"tags" binding:"required"`
	Files     []ArticlesFilesDBModel `json:"files" binding:"required"`
	CreatedAt time.Time              `json:"created_at" binding:"required"`
	UpdatedAt time.Time              `json:"updated_at" binding:"required"`
}

/* Data model for article index */
type ArticlesFilesIndexModel struct {
	Index int `json:"index" binding:"required" db:"index"`
}
