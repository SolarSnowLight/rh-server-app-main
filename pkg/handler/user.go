package handler

import (
	"encoding/json"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary GetProfile
// @Tags profile
// @Description Получение информации о профиле
// @ID user-profile-get
// @Accept  json
// @Produce  json
// @Param input body userModel.UserProfileDataModel true "credentials"
// @Success 200 {object} userModel.UserProfileDataModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/profile/get [post]
func (h *Handler) getProfile(c *gin.Context) {
	data, err := h.services.User.GetProfile(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var userProfile userModel.UserProfileDataModel

	err = json.Unmarshal([]byte(data.Data), &userProfile)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userProfile.Email = data.Email

	c.JSON(http.StatusOK, userProfile)
}

// @Summary UpdateProfile
// @Tags profile
// @Description Обновление информации о пользователе
// @ID user-profile-update
// @Accept  json
// @Produce  json
// @Param input body userModel.UserProfileUpdateDataModel true "credentials"
// @Success 200 {object} userModel.UserJSONBModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/profile/update [post]
func (h *Handler) updateProfile(c *gin.Context) {
	var input userModel.UserProfileUpdateDataModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.User.UpdateProfile(c, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary GetUserCompany
// @Tags profile
// @Description Получение информации о компании, к которой принадлежит пользователь
// @ID user-company-get
// @Accept  json
// @Produce  json
// @Success 200 {object} company.CompanyDbModelEx "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/company/get [post]
func (h *Handler) getUserCompany(c *gin.Context) {
	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.User.GetUserCompany(userId, domainId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(data.Uuid) <= 0 {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// @Summary GetUserRoles
// @Tags profile
// @Description Получение информации о всех ролях пользователя
// @ID user-role-get-all
// @Accept  json
// @Produce  json
// @Success 200 {object} company.CompanyDbModelEx "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/role/get/all [post]
func (h *Handler) getUserRoles(c *gin.Context) {
	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.User.GetAllRoles(userModel.UserIdentityModel{
		UserId:   userId,
		DomainId: domainId,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary CheckAccess
// @Tags profile
// @Description Проверка пользовательских прав на подключение к тому или иному административному модулю
// @ID user-check-access
// @Accept  json
// @Produce  json
// @Success 200 {object} BooleanResponse "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/access/check [post]
func (h *Handler) accessCheck(c *gin.Context) {
	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	var input rbacModel.RoleValueModel
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.services.User.AccessCheck(userId, domainId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, BooleanResponse{
		Value: data,
	})
}
