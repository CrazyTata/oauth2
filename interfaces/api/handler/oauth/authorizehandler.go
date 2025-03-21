package oauth

import (
	"log"
	"net/http"
	"oauth2/infrastructure/svc"

	"github.com/openshift/osin"
)

// AuthorizeHandler 处理获取用户课程的请求
func AuthorizeHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := newOAuthServer(svc)
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			// 验证客户端
			if ar.Client == nil {
				resp.SetError("unauthorized_client", "客户端未授权")
				osin.OutputJSON(resp, w, r)
				return
			}

			// 验证重定向URI
			if ar.RedirectUri == "" {
				resp.SetError("invalid_request", "缺少重定向URI")
				osin.OutputJSON(resp, w, r)
				return
			}

			// 设置授权
			ar.Authorized = true
			log.Printf("授权请求配置: Server.Config.AccessExpiration=%d", server.Config.AccessExpiration)

			// 完成授权请求,这里只会返回授权码
			server.FinishAuthorizeRequest(resp, r, ar)

			// 如果没有错误,会重定向到客户端的redirect_uri,并带上授权码
			if !resp.IsError {
				resp.Type = osin.REDIRECT
			}
		}

		// 输出响应(可能是重定向或错误信息)
		osin.OutputJSON(resp, w, r)
	}
}
