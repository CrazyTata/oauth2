package main

import (
	"flag"
	"fmt"
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
	api.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
