package route

const (
	AUTH_MAIN_ROUTE = "/auth"
)

const (
	// LOCAL
	AUTH_SIGN_IN_ROUTE        = "/sign-in"
	AUTH_SIGN_UP_ROUTE        = "/sign-up"
	AUTH_UPLOAD_PROFILE_IMAGE = "/sign-up/upload/image"

	// VK
	AUTH_SIGN_IN_VK_ROUTE          = "/sign-in/vk"
	AUTH_SIGN_IN_VK_CALLBACK_ROUTE = "/sign-in/vk/callback"

	// Google
	AUTH_SIGN_IN_GOOGLE_ROUTE = "/sign-in/oauth2"

	// MAIN
	AUTH_REFRESH_TOKEN_ROUTE = "/refresh"
	AUTH_LOGOUT_ROUTE        = "/logout"
	AUTH_ACTIVATE_ROUTE      = "/activate/:link"

	// Password Recovery
	AUTH_RECOVERY_PASSWORD = "/recovery/password"
	AUTH_RESET_PASSWORD    = "/reset/password"
)
