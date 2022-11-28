package service

import (
	"main-server/pkg/model/admin"
	adminModel "main-server/pkg/model/admin"
	repository "main-server/pkg/repository"

	"github.com/gin-gonic/gin"
)

/* Structure for this service */
type AdminService struct {
	repo repository.Admin
}

/* Function for create new struct of AdminService */
func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

/* Method for get all users, when location in system */
func (s *AdminService) GetAllUsers(c *gin.Context) (admin.UsersResponseModel, error) {
	return s.repo.GetAllUsers(c)
}

/* Method for create new company */
func (s *AdminService) CreateCompany(c *gin.Context, data adminModel.CompanyModel) (adminModel.CompanyModel, error) {
	return s.repo.CreateCompany(c, data)
}
