package oauth

import (
	"encoding/json"
	"net/http"
	"oauth2/infrastructure/svc"

	"github.com/openshift/osin"
)

// RefreshTokenHandler 处理刷新token的请求
func RefreshTokenHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 解析请求
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "无法解析请求参数",
			})
			return
		}

		// 获取refresh_token
		refreshToken := r.Form.Get("refresh_token")
		if refreshToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "缺少refresh_token参数",
			})
			return
		}

		// 初始化 OAuth 服务器
		server := newOAuthServer(svc)

		// 加载refresh token对应的访问数据
		accessData, err := server.Storage.LoadRefresh(refreshToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "无效的refresh_token",
			})
			return
		}

		// 创建新的访问令牌请求
		ar := &osin.AccessRequest{
			Type:            osin.REFRESH_TOKEN,
			Code:            "",
			Client:          accessData.Client,
			RedirectUri:     accessData.RedirectUri,
			Scope:           r.Form.Get("scope"),
			GenerateRefresh: true,
			Authorized:      true,
			Expiration:      server.Config.AccessExpiration,
		}

		// 处理访问令牌请求
		resp := server.NewResponse()
		defer resp.Close()

		// 移除旧的访问令牌
		if err := server.Storage.RemoveAccess(accessData.AccessToken); err != nil {
			resp.SetError("server_error", "无法移除旧的访问令牌")
			osin.OutputJSON(resp, w, r)
			return
		}

		// 移除旧的刷新令牌
		if err := server.Storage.RemoveRefresh(refreshToken); err != nil {
			resp.SetError("server_error", "无法移除旧的刷新令牌")
			osin.OutputJSON(resp, w, r)
			return
		}

		// 生成新的访问令牌
		server.FinishAccessRequest(resp, r, ar)
		if resp.IsError {
			osin.OutputJSON(resp, w, r)
			return
		}

		// 返回新的令牌
		osin.OutputJSON(resp, w, r)
	}
}
