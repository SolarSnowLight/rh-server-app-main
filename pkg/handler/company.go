package handler

import (
	pathConstant "main-server/pkg/constant/path"
	companyModel "main-server/pkg/model/company"
	userModel "main-server/pkg/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Summary GetManagers
// @Tags company
// @Description Получение среза из общего числа менеджеров компании
// @ID company-manager-get-all
// @Accept  json
// @Produce  json
// @Param input body companyModel.ManagerCountModel true "credentials"
// @Success 200 {object} companyModel.ManagerAnyCountModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/manager/get/all [post]
func (h *Handler) getManagers(c *gin.Context) {
	var input companyModel.ManagerCountModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Company.GetManagers(userId, domainId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary companyUpdateImage
// @Tags company
// @Description Обновление изображения компании
// @ID company-update-image
// @Accept  json
// @Produce  json
// @Param uuid query string true "uuid"
// @Param logo query string true "logo"
// @Success 200 {object} companyModel.CompanyImageModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/update [post]
func (h *Handler) companyUpdateImage(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	file := form.File["logo"]
	uuidCompany := c.PostForm("uuid")

	if (len(file) <= 0) || (uuidCompany == "") {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fileImg := file[len(file)-1]

	var data companyModel.CompanyImageModel

	newFilename := uuid.NewV4().String()
	filepath := pathConstant.PUBLIC_COMPANY + newFilename

	data, err = h.services.Company.CompanyUpdateImage(
		userModel.UserIdentityModel{
			UserId:   userId,
			DomainId: domainId,
		},
		companyModel.CompanyImageModel{
			Uuid:     uuidCompany,
			Filepath: filepath,
		},
	)

	if err != nil {
		form.RemoveAll()
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SaveUploadedFile(fileImg, filepath)
	c.JSON(http.StatusOK, data)
}

// @Summary UpdateCompany
// @Tags company
// @Description Создание новой компании (доступно только администратору)
// @ID company-update
// @Accept  json
// @Produce  json
// @Success 200 {object} companyModel.CompanyUpdateModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/update [post]
func (h *Handler) companyUpdate(c *gin.Context) {
	var input companyModel.CompanyUpdateModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Company.CompanyUpdate(
		userModel.UserIdentityModel{
			UserId:   userId,
			DomainId: domainId,
		},
		input,
	)

	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary Get manager
// @Tags company
// @Description Получение информации о менеджере
// @ID company-get-manager
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен доступа для текущего пользователя" example(Bearer access_token)
// @Param input body companyModel.ManagerUuidModel true "credentials"
// @Success 200 {object} companyModel.ManagerCompanyModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/manager/get [post]
func (h *Handler) companyGetManager(c *gin.Context) {
	var input companyModel.ManagerUuidModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Company.GetManager(
		userModel.UserIdentityModel{
			UserId:   userId,
			DomainId: domainId,
		},
		input,
	)

	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
