package service

import (
	rbacModel "main-server/pkg/model/rbac"
	repository "main-server/pkg/repository"
)

/* Structure for this service */
type DomainService struct {
	repo repository.Domain
}

/* Function for create new service */
func NewDomainService(repo repository.Domain) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

/* Get all articles */
func (s *DomainService) GetDomain(column, value interface{}) (rbacModel.DomainModel, error) {
	return s.repo.GetDomain(column, value)
}
