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
	docker image rm purple

docker_restart: docker_down docker_up

docker_purge_restart: docker_remove docker_up

local: swag
	docker compose up pg search-engine -d --build
	cd go-backend && go run cmd/main.go

proto_py:
	python -m grpc_tools.protoc -Iproto --python_out=python-backend --pyi_out=python-backend --grpc_python_out=python-backend proto/search_engine.proto

proto_go:
	protoc --go_out=go-backend/proto --go_opt=Mproto/search_engine.proto=. \
		--go-grpc_out=go-backend/proto --go-grpc_opt=Mproto/search_engine.proto=. \
		proto/search_engine.proto