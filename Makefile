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

model:
	#@goctl model mysql ddl -src infrastructure/sql/user.sql -dir infrastructure/persistence/model/user -c
	#@goctl model mysql ddl -src infrastructure/sql/course.sql -dir infrastructure/persistence/model/course -c
	#@goctl model mysql ddl -src infrastructure/sql/course_chapters.sql -dir infrastructure/persistence/model/course -c
	#@goctl model mysql ddl -src infrastructure/sql/my_course.sql -dir infrastructure/persistence/model/user -c
	#@goctl model mysql ddl -src infrastructure/sql/my_favorite.sql -dir infrastructure/persistence/model/user -c
	#@goctl model mysql ddl -src infrastructure/sql/notifications.sql -dir infrastructure/persistence/model/user -c
	#@goctl model mysql ddl -src infrastructure/sql/learning_record.sql -dir infrastructure/persistence/model/learning -c
	#@goctl model mysql ddl -src infrastructure/sql/learning_record_detail.sql -dir infrastructure/persistence/model/learning -c
	#@goctl model mysql ddl -src infrastructure/sql/record.sql -dir infrastructure/persistence/model/record -c
	#@goctl model mysql ddl -src infrastructure/sql/orders.sql -dir infrastructure/persistence/model/orders -c
	#@goctl model mysql ddl -src infrastructure/sql/apple_pay_subscription.sql -dir infrastructure/persistence/model/orders -c
	@goctl model mysql ddl -src infrastructure/sql/daily_metrics.sql -dir infrastructure/persistence/model/statistics -c