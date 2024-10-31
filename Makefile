migrations_dir = migrations

build:
	go build ./cmd/wikigraph

migration:
	goose -dir $(migrations_dir) create $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name') sql
