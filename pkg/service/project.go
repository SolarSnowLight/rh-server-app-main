package service

import (
	projectModel "main-server/pkg/model/project"
	userModel "main-server/pkg/model/user"
	repository "main-server/pkg/repository"
)

/* Structure for this service */
type ProjectService struct {
	repo repository.Project
}

/* Function for create new struct of ProjectService */
func NewProjectService(repo repository.Project) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

/* Method for create new project */
func (s *ProjectService) CreateProject(userId, domainId int, data projectModel.ProjectModel) (projectModel.ProjectModel, error) {
	return s.repo.CreateProject(userId, domainId, data)
}

/* Method for update project */
func (s *ProjectService) ProjectUpdate(user userModel.UserIdentityModel, data projectModel.ProjectUpdateModel) (projectModel.ProjectUpdateModel, error) {
	return s.repo.ProjectUpdate(user, data)
}

/* Method for update image for project */
func (s *ProjectService) ProjectUpdateImage(userId, domainId int, data projectModel.ProjectImageModel) (projectModel.ProjectImageModel, error) {
	return s.repo.ProjectUpdateImage(userId, domainId, data)
}

/* Method for get information about object */
func (s *ProjectService) GetProject(userId, domainId int, data projectModel.ProjectUuidModel) (projectModel.ProjectDbModel, error) {
	return s.repo.GetProject(userId, domainId, data)
}

/* Method for get any count projects */
func (s *ProjectService) GetProjects(userId, domainId int, data projectModel.ProjectCountModel) (projectModel.ProjectAnyCountModel, error) {
	return s.repo.GetProjects(userId, domainId, data)
}
