.PHONY: swag docker_up docker_down docker_remove docker_restart docker_purge_restart local

swag:
	~/go/bin/swag init -g go-backend/cmd/main.go -o go-backend/docs
	cd go-backend && ~/go/bin/swag fmt

docker_up:
	docker compose up -d --build

docker_down:
	docker compose down

docker_remove: docker_down
	docker volume rm purple_pg_data
	docker volume rm purple_redis_data
	docker image rm purple

docker_restart: docker_down docker_up

docker_purge_restart: docker_remove docker_up

local: swag
	docker compose up pg redis -d
	cd go-backend && go run cmd/main.go