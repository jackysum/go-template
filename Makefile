app_name = CHANGEME
dev_db = $(app_name)_dev
test_db = $(app_name)_test

createdb:
	createdb $(dev_db) && createdb $(test_db)

dropdb:
	dropdb $(dev_db) && dropdb $(test_db)

migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir migrations $$name

migrate:
	migrate -source file://migrations -database "postgres://localhost:5432/$(dev_db)?sslmode=disable" up && \
		migrate -source file://migrations -database "postgres://localhost:5432/$(test_db)?sslmode=disable" up

rollback:
	migrate -source file://migrations -database "postgres://localhost:5432/$(dev_db)?sslmode=disable" down 1 && \
		migrate -source file://migrations -database "postgres://localhost:5432/$(test_db)?sslmode=disable" down 1
