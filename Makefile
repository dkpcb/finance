dev/run/import:
	docker-compose up -d --build
	docker-compose exec server go run main.go import

dev/run/server:
	docker-compose up -d --build
	docker-compose exec server go run main.go