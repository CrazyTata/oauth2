package main

import (
	"flag"
	"fmt"
	"net/http"
	"oauth2/infrastructure/svc"

	"oauth2/common/redis"
	"oauth2/infrastructure/config"
	"oauth2/interfaces/api"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/config.yaml", "配置文件路径")

func main() {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	redis.Init(c.Redis.Host, c.Redis.Pass)
	defer redis.Close()

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 配置 CORS
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 允许的前端域名
			allowOrigins := []string{"http://127.0.0.1:8883", "http://localhost:8883", "http://localhost:3000", c.Domain}
			origin := r.Header.Get("Origin")
			for _, allowOrigin := range allowOrigins {
				if origin == allowOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	})

	api.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
