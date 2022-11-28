package handler

import (
	config "main-server/config"
	middlewareConstant "main-server/pkg/constant/middleware"
	pathConstant "main-server/pkg/constant/path"
	userModel "main-server/pkg/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// @Summary SignUp
// @Tags auth
// @Description Регистрация нового пользователя
// @ID auth-sign-up
// @Accept  json
// @Produce  json
// @Param input body userModel.UserRegisterModel true "account info"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input userModel.UserRegisterModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	c.SetSameSite(config.HTTPSameSite)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary UploadProfileImage
// @Tags auth
// @Description Загрузка пользовательского изображения (аватар)
// @ID auth-sign-up-upload-image
// @Accept  json
// @Produce  json
// @Param input body userModel.UserRegisterModel true "account info"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up/upload/image [post]
func (h *Handler) uploadProfileImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get data from form
	images := form.File["file"]
	profileImage := images[len(images)-1]

	filepath := pathConstant.PUBLIC_USER + uuid.NewV4().String()

	// Save profile image
	c.SaveUploadedFile(profileImage, filepath)

	_, err = h.services.Authorization.UploadProfileImage(c, filepath)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Изображение профиля пользователя было обновлено",
	})
}

// @Summary SignIn
// @Tags auth
// @Description Авторизация пользователя
// @ID auth-sign-in
// @Accept  json
// @Produce  json
// @Param input body userModel.UserLoginModel true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input userModel.UserLoginModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.Authorization.LoginUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	c.SetSameSite(config.HTTPSameSite)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary SignInVK
// @Tags auth
// @Description Авторизация пользователя через VK
// @ID auth-sign-in-vk
// @Accept  json
// @Produce  json
// @Param input body userModel.UserLoginModel true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in/vk [post]
func (h *Handler) signInVK(c *gin.Context) {
	var input userModel.UserLoginModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.Authorization.LoginUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	c.SetSameSite(config.HTTPSameSite)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary SignInOAuth2
// @Tags auth
// @Description Авторизация пользователя через Google OAuth2
// @ID auth-sign-in-oauth2
// @Accept  json
// @Produce  json
// @Param input body userModel.GoogleOAuth2Code true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in/oauth2 [post]
func (h *Handler) signInOAuth2(c *gin.Context) {
	var input userModel.GoogleOAuth2Code

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// For fast tests
	/*token, _ := configs.AppOAuth2Config.GoogleLogin.Exchange(c, input.Code)
	_, _ = google_oauth2.RevokeToken(token.AccessToken)
	return*/

	data, err := h.services.Authorization.LoginUserOAuth2(input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	c.SetSameSite(config.HTTPSameSite)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary refresh
// @Tags auth
// @Description Обновление токена доступа и токена обновления
// @ID auth-refresh
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен доступа для текущего пользователя" example(Bearer access_token)
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {

	// Получение токена обновления из файла cookie
	refreshToken, err := c.Cookie(viper.GetString("environment.refresh_token_key"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Получение дополнительной авторизационной информации пользователя (после работы middleware цепочки)
	accessToken, _ := c.Get(middlewareConstant.ACCESS_TOKEN_CTX)
	authTypeValue, _ := c.Get(middlewareConstant.AUTH_TYPE_VALUE_CTX)
	tokenApi, _ := c.Get(middlewareConstant.TOKEN_API_CTX)

	// Обновление токена доступа
	data, err := h.services.Authorization.Refresh(userModel.TokenLogoutDataModel{
		AccessToken:   accessToken.(string),
		RefreshToken:  refreshToken,
		AuthTypeValue: authTypeValue.(string),
		TokenApi:      tokenApi.(*string),
	}, refreshToken)

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Установка нового токена обновления (необходимо, если токен обновления изменился)
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	c.SetSameSite(config.HTTPSameSite)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

type LogoutOutputModel struct {
	IsLogout bool `json:"is_logout"`
}

// @Summary Logout
// @Tags auth
// @Description Выход из аккаунта
// @ID auth-logout
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен доступа для текущего пользователя" example(Bearer access_token)
// @Success 200 {object} LogoutOutputModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie(viper.GetString("environment.refresh_token_key"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, _ := c.Get(middlewareConstant.ACCESS_TOKEN_CTX)
	authTypeValue, _ := c.Get(middlewareConstant.AUTH_TYPE_VALUE_CTX)
	tokenApi, _ := c.Get(middlewareConstant.TOKEN_API_CTX)

	data, err := h.services.Authorization.Logout(userModel.TokenLogoutDataModel{
		AccessToken:   accessToken.(string),
		RefreshToken:  refreshToken,
		AuthTypeValue: authTypeValue.(string),
		TokenApi:      tokenApi.(*string),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if data {
		c.SetCookie(viper.GetString("environment.refresh_token_key"), "",
			30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
		c.SetSameSite(config.HTTPSameSite)
	}

	c.JSON(http.StatusOK, LogoutOutputModel{
		IsLogout: data,
	})
}

// @Summary Activate
// @Tags auth
// @Description Активация аккаунта по почте
// @ID auth-activate
// @Accept  json
// @Produce  json
// @Success 200 {object} LogoutOutputModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/activate [get]
func (h *Handler) activate(c *gin.Context) {
	_, err := h.services.Activate(c.Params.ByName("link"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "account_activate.html", gin.H{
		"title": "Подтверждение аккаунта",
	})
}

// @Summary Recovery password
// @Tags auth
// @Description Запрос на смену пароля пользователем
// @ID auth-recovery-password
// @Accept  json
// @Produce  json
// @Param input body userModel.UserEmailModel true "credentials"
// @Success 200 {object} successResponse "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/recovery/password [post]
func (h *Handler) recoveryPassword(c *gin.Context) {
	var input userModel.UserEmailModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	_, err := h.services.Authorization.RecoveryPassword(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Message: "На Вашу почту была отправлена ссылка с подтверждением изменения пароля",
	})
}

// @Summary Reset password
// @Tags auth
// @Description Изменение пароля пользователем
// @ID auth-reset-password
// @Accept  json
// @Produce  json
// @Param input body userModel.ResetPasswordModel true "credentials"
// @Success 200 {object} successResponse "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/reset/password [post]
func (h *Handler) resetPassword(c *gin.Context) {
	var input userModel.ResetPasswordModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	_, err := h.services.Authorization.ResetPassword(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Message: "Пароль был успешно изменён!",
	})
}
