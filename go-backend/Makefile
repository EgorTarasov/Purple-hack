.PHONY: swag docker_up docker_down docker_remove docker_restart docker_purge_restart local

swag:
	~/go/bin/swag init -g cmd/main.go -o ./docs
	~/go/bin/swag fmt

docker_up:
	docker compose up -d --build

docker_down:
	docker compose down

docker_remove: docker_down
	docker volume rm hack-template_pg_data
	docker volume rm hack-template_redis_data
	docker image rm hack-template

docker_restart: docker_down docker_up

docker_purge_restart: docker_remove docker_up

local: swag
	docker compose up pg redis -d
	go run cmd/main.go