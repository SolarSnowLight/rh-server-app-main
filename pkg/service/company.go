package service

import (
	companyModel "main-server/pkg/model/company"
	userModel "main-server/pkg/model/user"
	repository "main-server/pkg/repository"
)

/* Structure for this service */
type CompanyService struct {
	repo repository.Company
}

/* Function for create new struct of CompanyService */
func NewCompanyService(repo repository.Company) *CompanyService {
	return &CompanyService{
		repo: repo,
	}
}

/* Method for create new project */
func (s *CompanyService) GetManagers(userId, domainId int, data companyModel.ManagerCountModel) (companyModel.ManagerAnyCountModel, error) {
	return s.repo.GetManagers(userId, domainId, data)
}

/* Method for update image of company */
func (s *CompanyService) CompanyUpdateImage(user userModel.UserIdentityModel, data companyModel.CompanyImageModel) (companyModel.CompanyImageModel, error) {
	return s.repo.CompanyUpdateImage(user, data)
}

/* Method for update company info */
func (s *CompanyService) CompanyUpdate(user userModel.UserIdentityModel, data companyModel.CompanyUpdateModel) (companyModel.CompanyUpdateModel, error) {
	return s.repo.CompanyUpdate(user, data)
}

/* Method for get information about define manager */
func (s *CompanyService) GetManager(user userModel.UserIdentityModel, data companyModel.ManagerUuidModel) (companyModel.ManagerCompanyModel, error) {
	return s.repo.GetManager(user, data)
}
