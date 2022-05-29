.PHONY: build
build:
	go build -v ./cmd/service

.PHONY: wait
wait:
	./wait-for-postgres.sh db

.PHONY: migrate_up
migrate_up:
	migrate -database "postgres://postgres:HEYO@db/postgres?sslmode=disable" -path migrations up

.PHONY: chmodfile
chmodfile:
	chmod +x ./wait-for-postgres.sh

.PHONY: vendor
chmodfile:
	go mod vendor

.PHONY: run
run: chmodfile wait vendor migrate_up build
	./service

.DEFAULT_GOAL := build