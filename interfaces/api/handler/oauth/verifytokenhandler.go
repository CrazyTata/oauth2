package oauth

import (
	"net/http"
	"oauth2/infrastructure/svc"
	"time"

	"github.com/openshift/osin"
)

// VerifyTokenHandler 处理验证access token的请求
func VerifyTokenHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 初始化 OAuth 服务器
		server := newOAuthServer(svc)
		resp := server.NewResponse()
		defer resp.Close()

		// 从请求头获取访问令牌
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			resp.SetError("invalid_token", "无效的Authorization头部")
			resp.StatusCode = http.StatusUnauthorized
			osin.OutputJSON(resp, w, r)
			return
		}
		accessToken := authHeader[7:]

		// 加载访问令牌
		accessData, err := server.Storage.LoadAccess(accessToken)
		if err != nil {
			resp.SetError("invalid_token", "访问令牌无效或已过期")
			resp.StatusCode = http.StatusUnauthorized
			osin.OutputJSON(resp, w, r)
			return
		}

		// 计算剩余有效期
		expiresIn := int32(time.Until(accessData.CreatedAt.Add(time.Duration(accessData.ExpiresIn) * time.Second)).Seconds())
		if expiresIn <= 0 {
			resp.SetError("invalid_token", "访问令牌已过期")
			resp.StatusCode = http.StatusUnauthorized
			osin.OutputJSON(resp, w, r)
			return
		}

		// 设置响应
		resp.Output = map[string]interface{}{
			"valid":       true,
			"client_id":   accessData.Client.GetId(),
			"scope":       accessData.Scope,
			"expires_in":  expiresIn,
			"create_time": accessData.CreatedAt.Format(time.RFC3339),
		}

		osin.OutputJSON(resp, w, r)
	}
}
