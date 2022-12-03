package service

import (
	companyModel "main-server/pkg/model/company"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	repository "main-server/pkg/repository"

	"github.com/gin-gonic/gin"
)

/* Structure for this service */
type UserService struct {
	repo repository.User
}

/* Function for create new service */
func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

/* Get information about profile user */
func (s *UserService) GetProfile(c *gin.Context) (userModel.UserProfileModel, error) {
	return s.repo.GetProfile(c)
}

/* Method for update profile user */
func (s *UserService) UpdateProfile(c *gin.Context, data userModel.UserProfileUpdateDataModel) (userModel.UserJSONBModel, error) {
	return s.repo.UpdateProfile(c, data)
}

/* A method for obtaining information about the company in which this user is present */
func (s *UserService) GetUserCompany(userId, domainId int) (companyModel.CompanyDbModelEx, error) {
	return s.repo.GetUserCompany(userId, domainId)
}

/* Method for check access */
func (s *UserService) AccessCheck(userId, domainId int, value rbacModel.RoleValueModel) (bool, error) {
	return s.repo.AccessCheck(userId, domainId, value)
}

func (s *UserService) GetAllRoles(user userModel.UserIdentityModel) (*userModel.UserRoleModel, error) {
	return s.repo.GetAllRoles(user)
}
