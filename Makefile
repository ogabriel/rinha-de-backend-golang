ifneq (,$(wildcard ./.env))
    include .env
	export
endif

PROJECT=$(shell basename $(PWD))
DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)

run:
	go run .

run-release:
	GIN_MODE=release go run .

database-reset:
	make database-drop || exit 0
	make database-create
	make database-migration-up

database-reset-release:
	make database-check
	make database-drop || exit 0
	make database-create
	/app/migrate -path migrations/ -database $(DATABASE_URL)/$(DATABASE_NAME)?sslmode=disable -verbose up

database-create:
	psql $(DATABASE_URL) -c "CREATE DATABASE $(DATABASE_NAME)"

database-drop:
	psql $(DATABASE_URL) -c "DROP DATABASE $(DATABASE_NAME)"

database-migration-up:
	migrate -path migrations/ -database $(DATABASE_URL)/$(DATABASE_NAME)?sslmode=disable -verbose up

database-migration-create:
	migrate create -ext sql -dir migrations -seq $(name)

database-check:
	until nc -z -v -w30 $(DATABASE_HOST) $(DATABASE_PORT); do \
	  sleep 1; \
	done

docker-compose-one:
	make docker-compose-down
	docker compose -f docker-compose.one.yml -p $(PROJECT)-one up

docker-compose-one-build:
	make docker-compose-down
	docker compose -f docker-compose.one.yml -p $(PROJECT)-one up --build

docker-compose-two:
	make docker-compose-down
	docker compose -f docker-compose.two.yml -p $(PROJECT)-two up

docker-compose-two-build:
	make docker-compose-down
	docker compose -f docker-compose.two.yml -p $(PROJECT)-two up --build

docker-compose-down:
	docker stop postgres-15 || exit 0
	docker stop postgres-11 || exit 0
	docker stop postgres || exit 0
	docker compose -p $(PROJECT)-dev down || exit 0
	docker compose -p $(PROJECT)-two down || exit 0
	docker compose -p $(PROJECT)-one down || exit 0
