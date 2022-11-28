package service

import (
	userModel "main-server/pkg/model/user"
	repository "main-server/pkg/repository"
)

// Структура репозитория
type AuthTypeService struct {
	authType repository.AuthType
}

// Функция создания нового репозитория
func NewAuthTypeService(role repository.AuthType) *AuthTypeService {
	return &AuthTypeService{authType: role}
}

func (s *AuthTypeService) GetAuthType(column, value string) (userModel.AuthTypeModel, error) {
	return s.authType.GetAuthType(column, value)
}
