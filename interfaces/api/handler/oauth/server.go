package oauth

import (
	"oauth2/application/service"
	"oauth2/infrastructure/svc"

	"github.com/openshift/osin"
)

// newOAuthServer 创建一个新的OAuth服务器实例
func newOAuthServer(svc *svc.ServiceContext) *osin.Server {
	config := osin.NewServerConfig()
	config.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE}
	config.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
	}
	config.AuthorizationExpiration = 600 // 10分钟
	config.AccessExpiration = 3600       // 1小时
	config.AllowGetAccessRequest = true
	config.ErrorStatusCode = 401

	storage := service.NewStorage(svc, "osin_")
	server := osin.NewServer(config, storage)

	return server
}
