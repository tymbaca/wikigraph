migrations_dir = migrations

migration:
	goose -dir $(migrations_dir) create $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name') sql