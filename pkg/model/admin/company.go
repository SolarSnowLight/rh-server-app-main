package admin

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
