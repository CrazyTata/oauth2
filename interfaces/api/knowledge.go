package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"oauth2/infrastructure/svc"
	"oauth2/interfaces/api/handler/oauth"
)

// RegisterHandlers 注册HTTP处理器
func RegisterHandlers(server *rest.Server, svc *svc.ServiceContext) {
	// 工具
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/oauth/authorize",
				Handler: oauth.AuthorizeHandler(svc),
			},
			{
				Method:  http.MethodPost,
				Path:    "/v1/oauth/token",
				Handler: oauth.TokenHandler(svc),
			},
			{
				Method:  http.MethodPost,
				Path:    "/v1/oauth/create-client",
				Handler: oauth.CreateClientHandler(svc),
			},
			{
				Method:  http.MethodGet,
				Path:    "/v1/oauth/callback",
				Handler: oauth.CallbackHandler(svc),
			},
		},
	)
}
