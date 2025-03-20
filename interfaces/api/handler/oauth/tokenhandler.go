package oauth

import (
	"net/http"
	"oauth2/infrastructure/svc"

	"github.com/openshift/osin"
	"github.com/zeromicro/go-zero/core/logx"
)

// TokenHandler 处理获取访问令牌的请求
func TokenHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logx.WithContext(r.Context())
		server := newOAuthServer(svc)
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			// 验证客户端
			if ar.Client == nil {
				resp.SetError("unauthorized_client", "客户端未授权")
				osin.OutputJSON(resp, w, r)
				return
			}

			// 根据不同的授权类型进行处理
			switch ar.Type {
			case osin.AUTHORIZATION_CODE:
				// 验证授权码
				if ar.AuthorizeData == nil {
					resp.SetError("invalid_grant", "授权码无效或已过期")
					osin.OutputJSON(resp, w, r)
					return
				}
			case osin.REFRESH_TOKEN:
				// 验证刷新令牌
				if r.FormValue("refresh_token") == "" {
					resp.SetError("invalid_grant", "刷新令牌无效")
					osin.OutputJSON(resp, w, r)
					return
				}
			case osin.CLIENT_CREDENTIALS:
				// 验证客户端凭证
				if ar.Client.GetSecret() == "" {
					resp.SetError("invalid_client", "客户端密钥无效")
					osin.OutputJSON(resp, w, r)
					return
				}
			default:
				resp.SetError("unsupported_grant_type", "不支持的授权类型")
				osin.OutputJSON(resp, w, r)
				return
			}

			// 授权请求
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}

		if resp.IsError {
			logger.Errorf("Token error: %v", resp.InternalError)
		} else {
			logger.Infof("Token granted: %s", resp.Output["access_token"])
		}

		osin.OutputJSON(resp, w, r)
	}
}
