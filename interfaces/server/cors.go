package server

const (
	AccessControlAllowOrigin      = "*"
	AccessControlAllowCredentials = "true"
	AccessControlAllowMethods     = "POST, OPTIONS, GET, PUT"
	AccessControlAllowHeaders     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
)

type CorsOptions struct {
	AccessControlAllowOrigin      string
	AccessControlAllowCredentials string
	AccessControlAllowHeaders     string
	AccessControlAllowMethods     string
}

func GetAllAllowedCors() CorsOptions {
	return CorsOptions{
		AccessControlAllowOrigin:      AccessControlAllowOrigin,
		AccessControlAllowHeaders:     AccessControlAllowHeaders,
		AccessControlAllowMethods:     AccessControlAllowMethods,
		AccessControlAllowCredentials: AccessControlAllowCredentials,
	}
}

func GetCorsOptions(
	AccessControlAllowOrigin,
	AccessControlAllowHeaders,
	AccessControlAllowMethods,
	AccessControlAllowCredentials string) CorsOptions {

	return CorsOptions{
		AccessControlAllowOrigin:      AccessControlAllowOrigin,
		AccessControlAllowHeaders:     AccessControlAllowHeaders,
		AccessControlAllowMethods:     AccessControlAllowMethods,
		AccessControlAllowCredentials: AccessControlAllowCredentials,
	}

}
