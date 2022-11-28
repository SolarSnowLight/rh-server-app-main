package handler

import (
	"errors"
	middlewareConstants "main-server/pkg/constant/middleware"
	authService "main-server/pkg/service/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(middlewareConstants.AUTHORIZATION_HEADER)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Пустой заголовок авторизации!")
		return
	}

	headerParts := strings.Split(header, " ")
	if (len(headerParts) != 2) || (headerParts[1] == "null") || (headerParts[1] == "undefined") {
		newErrorResponse(c, http.StatusUnauthorized, "Пользователь не авторизован!")
		return
	}

	data, err := h.services.Token.ParseToken(headerParts[1], viper.GetString("token.signing_key_access"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	domain, err := h.services.Domain.GetDomain("value", viper.GetString("domain"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	switch data.AuthType.Value {
	case "GOOGLE":
		if result, err := authService.VerifyAccessToken(*data.TokenApi); err != nil || result != true {
			newErrorResponse(c, http.StatusUnauthorized, "Не действительный токен доступа")
			return
		}
		break

	case "LOCAL":
		break
	}

	// Добавление к контексту дополнительных данных о пользователе
	c.Set(middlewareConstants.USER_CTX, data.UsersId)
	c.Set(middlewareConstants.AUTH_TYPE_VALUE_CTX, data.AuthType.Value)
	c.Set(middlewareConstants.TOKEN_API_CTX, data.TokenApi)
	c.Set(middlewareConstants.ACCESS_TOKEN_CTX, headerParts[1])
	c.Set(middlewareConstants.DOMAINS_ID, domain.Id)
}

func (h *Handler) userIdentityLogout(c *gin.Context) {
	header := c.GetHeader(middlewareConstants.AUTHORIZATION_HEADER)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Пустой заголовок авторизации!")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Не корректный авторизационный заголовок!")
		return
	}

	data, err := h.services.Token.ParseTokenWithoutValid(headerParts[1], viper.GetString("token.signing_key_access"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Добавление к контексту дополнительных данных о пользователе
	c.Set(middlewareConstants.USER_CTX, data.UsersId)
	c.Set(middlewareConstants.AUTH_TYPE_VALUE_CTX, data.AuthType.Value)
	c.Set(middlewareConstants.TOKEN_API_CTX, data.TokenApi)
	c.Set(middlewareConstants.ACCESS_TOKEN_CTX, headerParts[1])
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(middlewareConstants.USER_CTX)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

/* Function for get information user with gin.Context */
func getContextUserInfo(c *gin.Context) (int, int, error) {
	usersId, exist := c.Get(middlewareConstants.USER_CTX)

	if !exist {
		return -1, -1, errors.New("Нет доступа!")
	}

	domainsId, exist := c.Get(middlewareConstants.DOMAINS_ID)

	if !exist {
		return -1, -1, errors.New("Нет доступа!")
	}

	return usersId.(int), domainsId.(int), nil
}

func (h *Handler) userIdentityHasRoles(exp string, roles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		usersId, domainsId, err := getContextUserInfo(c)
		if err != nil {
			newErrorResponse(c, http.StatusForbidden, "Нет доступа!")
			return
		}

		flags := make([]bool, 0)

		for _, element := range roles {
			has, err := h.services.Role.HasRole(usersId, domainsId, element)

			if err != nil {
				newErrorResponse(c, http.StatusForbidden, "Нет доступа!")
				return
			}

			flags = append(flags, has)
		}

		access := false

		if exp == "AND" || exp == "and" {
			access = true
			for _, element := range flags {
				if !element {
					access = false
					break
				}
			}
		} else if exp == "OR" || exp == "or" {
			access = false
			for _, element := range flags {
				if element {
					access = true
					break
				}
			}
		}

		if !access {
			newErrorResponse(c, http.StatusForbidden, "Нет доступа!")
			return
		}
	}
}

func (h *Handler) userIdentityHasRole(role string) func(c *gin.Context) {
	return func(c *gin.Context) {
		usersId, domainsId, err := getContextUserInfo(c)
		if err != nil {
			newErrorResponse(c, http.StatusForbidden, "Нет доступа!")
			return
		}

		has, err := h.services.Role.HasRole(usersId, domainsId, role)

		if (err != nil) || (!has) {
			newErrorResponse(c, http.StatusForbidden, "Нет доступа!")
			return
		}
	}
}
