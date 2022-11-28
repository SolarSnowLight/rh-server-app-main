package service

import (
	rbacModel "main-server/pkg/model/rbac"
	repository "main-server/pkg/repository"
)

/* Structure for this service */
type RoleService struct {
	repo repository.Role
}

/* Function for create new service */
func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

/* Get role */
func (s *RoleService) GetRole(column, value interface{}) (rbacModel.RoleModel, error) {
	return s.repo.GetRole(column, value)
}

/* HasRole */
func (s *RoleService) HasRole(usersId, domainsId int, roleValue string) (bool, error) {
	return s.repo.HasRole(usersId, domainsId, roleValue)
}
