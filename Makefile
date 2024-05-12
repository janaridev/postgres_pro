include .env
export

CONN_LINK := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_EXTERNAL_PORT)/$(PG_DB_NAME)?sslmode=$(PG_USE_SSL)

migration_create:
	migrate create -ext sql -dir migrations/postgres -seq $(NAME)

migration_up:
	migrate -path migrations/postgres -database "$(CONN_LINK)" -verbose up

migration_down:
	migrate -path migrations/postgres -database "$(CONN_LINK)" -verbose down

migration_fix:
	migrate -path migrations/postgres/ -database "$(CONN_LINK)" force $(VERSION)

test:
	go test -v ./...

docker_build:
	docker build -f ./docker/Dockerfile -t csharpjanari/postgres-pro .