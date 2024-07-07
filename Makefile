# Docker Compose でサービスをビルドしてバックグラウンドで起動
dev/run/server:
	docker-compose up -d --build
	docker-compose exec server go run main.go