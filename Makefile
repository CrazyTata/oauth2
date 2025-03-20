.PHONY: run 
# 构建
build:
	@echo "拉取代码.."
	@git pull
	@echo "构建镜像..."
	@docker-compose build oauth2
	@echo "停止容器..."
	@docker-compose stop oauth2
	@echo "启动容器..."
	@docker-compose up -d oauth2
