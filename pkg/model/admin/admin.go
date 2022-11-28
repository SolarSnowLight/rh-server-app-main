package admin

/* Model data for response users model */
type UsersResponseModel struct {
	Users *[]UserResponseModel `json:"users" binding:"required"`
}

/* Model data for response user model */
type UserResponseModel struct {
	Email string `json:"email" binding:"required" db:"email"`
}
