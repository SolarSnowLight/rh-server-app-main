package handler

import (
	"fmt"
	roleConstant "main-server/pkg/constant/role"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "main-server/docs"

	route "main-server/pkg/constant/route"
	service "main-server/pkg/service"

	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

/* Инициализация маршрутов */
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Set max multipart memory
	router.MaxMultipartMemory = 50 << 20 // 50 MiB

	// Set folder for storage static files
	router.Static("/public", "./public")

	// Set global folder for load html
	router.LoadHTMLGlob("pkg/template/*")

	// Settings cors policies
	router.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins:     []string{viper.GetString("client_url")},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		AllowCredentials: true,
	}))

	// URL: /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// URL: /auth
	auth := router.Group(route.AUTH_MAIN_ROUTE)
	{
		// URL: /auth/sign-up
		auth.POST(route.AUTH_SIGN_UP_ROUTE, h.signUp)

		// URL: /auth/sign-in
		auth.POST(route.AUTH_SIGN_IN_ROUTE, h.signIn)

		// URL: /auth/sign-in/oauth2
		auth.POST(route.AUTH_SIGN_IN_GOOGLE_ROUTE, h.signInOAuth2)

		// URL: /auth/activate/:link
		auth.GET(route.AUTH_ACTIVATE_ROUTE, h.activate)

		// URL: /auth/refresh
		auth.POST(route.AUTH_REFRESH_TOKEN_ROUTE, h.userIdentityLogout, h.refresh)

		// URL: /auth/logout
		auth.POST(route.AUTH_LOGOUT_ROUTE, h.userIdentityLogout, h.logout)

		// URL: /auth/sign-up/upload/image
		auth.POST(route.AUTH_UPLOAD_PROFILE_IMAGE, h.userIdentity, h.uploadProfileImage)

		// URL: /auth/recovery/password
		auth.POST(route.AUTH_RECOVERY_PASSWORD, h.recoveryPassword)

		// URL: /auth/reset/password
		auth.POST(route.AUTH_RESET_PASSWORD, h.resetPassword)
	}

	// URL: /user
	user := router.Group(route.USER_MAIN_ROUTE, h.userIdentity)
	{
		// URL: /user/access/check
		user.POST(route.USER_CHECK_ACCESS_ROUTE, h.accessCheck)

		// URL: /user/profile
		profile := user.Group(route.USER_PROFILE_ROUTE)
		{
			// URL: /user/profile/get
			profile.POST(route.GET_ROUTE, h.getProfile)

			// URL: /user/profile/update
			profile.POST(route.UPDATE_ROUTE, h.updateProfile)
		}

		// URL: /user/company
		company := user.Group(route.COMPANY_MAIN_ROUTE)
		{
			// URL: /user/company/get
			company.POST(route.GET_ROUTE, h.getUserCompany)
		}
	}

	// URL: /admin
	admin := router.Group(route.ADMIN_MAIN_ROUTE, h.userIdentity, h.userIdentityHasRole(roleConstant.ROLE_ADMIN))
	{
		// URL: /admin/user
		user := admin.Group(route.ADMIN_USER)
		{
			// URL: /admin/user/get/all
			user.POST(route.GET_ALL_ROUTE, h.getAllUsers)
		}

		// URL: /admin/company
		company := admin.Group(route.ADMIN_COMPANY)
		{
			// URL: /admin/company/create
			company.POST(route.CREATE_ROUTE, h.createCompany)
		}
	}

	// URL: /company
	company := router.Group(
		route.COMPANY_MAIN_ROUTE,
		h.userIdentity,
		h.userIdentityHasRoles("OR", roleConstant.ROLE_ADMIN, roleConstant.ROLE_BUILDER_ADMIN),
	)
	{
		// URL: /project
		project := company.Group(route.PROJECT_MAIN_ROUTE)
		{
			// URL: /company/project/create
			project.POST(route.CREATE_ROUTE, h.userIdentityHasRole(roleConstant.ROLE_BUILDER_ADMIN), h.createProject)

			// URL: /company/project/update
			project.POST(
				route.UPDATE_ROUTE,
				h.userIdentityHasRoles(
					"OR",
					roleConstant.ROLE_BUILDER_MANAGER,
					roleConstant.ROLE_BUILDER_ADMIN,
				),
				h.projectUpdate,
			)

			// URL: /company/project/update/image
			project.POST(
				fmt.Sprintf("%s/%s", route.UPDATE_ROUTE, route.RESOURCE_IMAGE_ROUTE),
				h.userIdentityHasRoles(
					"OR",
					roleConstant.ROLE_BUILDER_MANAGER,
					roleConstant.ROLE_BUILDER_ADMIN,
				),
				h.projectUpdateImage,
			)

			// URL: /company/project/get
			project.POST(route.GET_ROUTE, h.getProject)

			// URL: /company/project/get/all
			project.POST(route.GET_ALL_ROUTE, h.getProjects)
		}

		// URL: /manager
		manager := company.Group(route.MANAGER_MAIN_ROUTE)
		{
			// URL: /company/manager/get/all
			manager.POST(route.GET_ALL_ROUTE, h.userIdentityHasRole(roleConstant.ROLE_BUILDER_ADMIN), h.getManagers)

			// URL: /company/manager/get
			manager.POST(route.GET_ROUTE, h.companyGetManager)
		}

		// URL: /company/create
		company.POST(route.CREATE_ROUTE, h.userIdentityHasRole(roleConstant.ROLE_ADMIN), h.createCompany)

		// URL: /company/update/image
		company.POST(fmt.Sprintf("%s/%s", route.UPDATE_ROUTE, route.RESOURCE_IMAGE_ROUTE),
			h.userIdentityHasRoles(
				"OR",
				roleConstant.ROLE_BUILDER_ADMIN,
				roleConstant.ROLE_ADMIN,
				roleConstant.ROLE_MANAGER,
				roleConstant.ROLE_SUPER_ADMIN,
			),
			h.companyUpdateImage,
		)

		// URL: /company/update
		company.POST(route.UPDATE_ROUTE,
			h.userIdentityHasRoles(
				"OR",
				roleConstant.ROLE_BUILDER_ADMIN,
				roleConstant.ROLE_ADMIN,
				roleConstant.ROLE_MANAGER,
				roleConstant.ROLE_SUPER_ADMIN,
			),
			h.companyUpdate,
		)
	}

	return router
}
