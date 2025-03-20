package oauth

import (
	"net/http"
	"oauth2/application/service"
	"oauth2/infrastructure/svc"
)

// CreateClientHandler 处理创建客户端的请求
func CreateClientHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := newOAuthServer(svc)
		storage := server.Storage.(*service.Storage)

		client := storage.CreateClientWithInformation(
			"1234",      // client_id
			"secret123", // client_secret
			"http://127.0.0.1:8883/v1/oauth/callback", // redirect_uri
			nil, // user_data
		)

		if err := storage.CreateClient(client); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Client created successfully"))
	}
}
