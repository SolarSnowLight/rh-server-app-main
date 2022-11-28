package company

/* Company UUID model */
type ManagerUuidModel struct {
	CompanyUuid string `json:"company_uuid" binding:"required"`
	ManagerUuid string `json:"manager_uuid" binding:"required"`
}

type CompanyImageModel struct {
	Uuid     string `json:"uuid" binding:"required"`
	Filepath string `json:"Filepath" binding:"required"`
}

/* Model data for company */
type CompanyModel struct {
	Logo         string `json:"logo" binding:"required" db:"logo"`
	Title        string `json:"title" binding:"required" db:"title"`
	Description  string `json:"description" binding:"required" db:"description"`
	Phone        string `json:"phone" binding:"required" db:"phone"`
	Link         string `json:"link" binding:"required" db:"link"`
	EmailCompany string `json:"email_company" binding:"required" db:"email_company"`
	EmailAdmin   string `json:"email_admin" binding:"required" db:"email_admin"`
}

type CompanyUpdateModel struct {
	Uuid         string  `json:"uuid" binding:"required" db:"uuid"`
	Title        string  `json:"title" binding:"required" db:"title"`
	Description  string  `json:"description" binding:"required" db:"description"`
	Phone        string  `json:"phone" binding:"required" db:"phone"`
	Link         string  `json:"link" binding:"required" db:"link"`
	EmailCompany string  `json:"email_company" binding:"required" db:"email_company"`
	EmailAdmin   *string `json:"email_admin" db:"email_admin"`
}

type CompanyDbModel struct {
	Uuid      string `json:"uuid" binding:"required" db:"uuid"`
	Data      string `json:"data" binding:"required" db:"data"`
	CreatedAt string `json:"created_at" binding:"required" db:"created_at"`
}

type CompanyDbModelEx struct {
	Uuid      string       `json:"uuid" binding:"required"`
	Data      CompanyModel `json:"data" binding:"required"`
	CreatedAt string       `json:"created_at" binding:"required"`
	Rules     []string     `json:"rules" binding:"required"`
}

type CompanyRuleModelEx struct {
	Value string `json:"value" binding:"required" db:"value"`
	V3    string `json:"v3" binding:"required" db:"v3"`
}
