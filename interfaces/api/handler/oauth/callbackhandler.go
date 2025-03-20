package oauth

import (
	"fmt"
	"net/http"
	"oauth2/infrastructure/svc"

	"github.com/openshift/osin"
)

// CallbackHandler 处理获取用户课程的请求
func CallbackHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取授权码
		code := r.URL.Query().Get("code")
		if code == "" {
			// 检查是否有错误信息
			if error := r.URL.Query().Get("error"); error != "" {
				errorDesc := r.URL.Query().Get("error_description")
				http.Error(w, fmt.Sprintf("授权失败: %s - %s", error, errorDesc), http.StatusBadRequest)
				return
			}
			http.Error(w, "未收到授权码", http.StatusBadRequest)
			return
		}

		// 初始化 OAuth 服务器
		server := newOAuthServer(svc)

		// 先加载授权数据
		authData, err := server.Storage.LoadAuthorize(code)
		if err != nil {
			resp := server.NewResponse()
			resp.SetError("invalid_grant", "授权码无效或已过期")
			osin.OutputJSON(resp, w, r)
			return
		}

		// 创建访问令牌请求
		ar := &osin.AccessRequest{
			Type:            osin.AUTHORIZATION_CODE,
			Code:            code,
			Client:          authData.Client,
			RedirectUri:     authData.RedirectUri,
			Scope:           authData.Scope,
			GenerateRefresh: true,
			Authorized:      true,
		}

		// 处理访问令牌请求
		resp := server.NewResponse()
		defer resp.Close()

		if err := server.Storage.RemoveAuthorize(code); err != nil {
			resp.SetError("server_error", "无法删除授权码")
			osin.OutputJSON(resp, w, r)
			return
		}

		server.FinishAccessRequest(resp, r, ar)
		if resp.IsError {
			osin.OutputJSON(resp, w, r)
			return
		}

		// API 请求则返回 JSON
		osin.OutputJSON(resp, w, r)
	}
}
