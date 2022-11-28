package handler

import (
	"encoding/json"
	pathConstant "main-server/pkg/constant/path"
	projectModel "main-server/pkg/model/project"
	userModel "main-server/pkg/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Summary CreateProject
// @Tags project
// @Description Создание нового проекта в компании
// @ID company-project-create
// @Accept  json
// @Produce  json
// @Param input body projectModel.ProjectModel true "credentials"
// @Success 200 {object} projectModel.ProjectModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/project/create [post]
func (h *Handler) createProject(c *gin.Context) {
	var input projectModel.ProjectModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Project.CreateProject(userId, domainId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary GetProject
// @Tags project
// @Description Получение информации о конкретном проекте
// @ID company-project-get
// @Accept  json
// @Produce  json
// @Param input body projectModel.ProjectUuidModel true "credentials"
// @Success 200 {object} projectModel.ProjectDbDataEx "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/project/get [post]
func (h *Handler) getProject(c *gin.Context) {
	var input projectModel.ProjectUuidModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Project.GetProject(userId, domainId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// From json to object
	var projectData projectModel.ProjectDataModel
	err = json.Unmarshal([]byte(data.Data), &projectData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, projectModel.ProjectDbDataEx{
		Uuid:      data.Uuid,
		Data:      projectData,
		CreatedAt: data.CreatedAt,
	})
}

// @Summary GetProjects
// @Tags project
// @Description Получение среза из общего числа проектов компании
// @ID company-project-get-all
// @Accept  json
// @Produce  json
// @Param input body projectModel.ProjectCountModel true "credentials"
// @Success 200 {object} projectModel.ProjectAnyCountModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/project/get/all [post]
func (h *Handler) getProjects(c *gin.Context) {
	var input projectModel.ProjectCountModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Project.GetProjects(userId, domainId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary ProjectUpdateImage
// @Tags project
// @Description Добавление нового логотипа проекта
// @ID company-project-add-logo
// @Accept  json
// @Produce  json
// @Param logo query string true "logo"
// @Param uuid query string true "uuid"
// @Success 200 {object} projectModel.ProjectImageModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/project/update/image [post]
func (h *Handler) projectUpdateImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file := form.File["logo"]
	fileImg := file[len(file)-1]
	projectUuid := c.PostForm("uuid")

	newFilename := uuid.NewV4().String()
	filepath := pathConstant.PUBLIC_PROJECT + newFilename

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Project.ProjectUpdateImage(
		userId,
		domainId,
		projectModel.ProjectImageModel{
			Filepath: filepath,
			Uuid:     projectUuid,
		},
	)

	if err != nil {
		form.RemoveAll()
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.SaveUploadedFile(fileImg, filepath)
	c.JSON(http.StatusOK, data)
}

// @Summary ProjectUpdate
// @Tags project
// @Description Обновление информации о проекте в компании
// @ID project-update
// @Accept  json
// @Produce  json
// @Success 200 {object} projectModel.ProjectUpdateModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /company/project/update [post]
func (h *Handler) projectUpdate(c *gin.Context) {
	var input projectModel.ProjectUpdateModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, domainId, err := getContextUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Project.ProjectUpdate(
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
