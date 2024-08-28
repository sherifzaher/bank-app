postgres:
	docker run --name postgres-test-1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest
createdb:
	docker exec -it postgres-test-1 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-test-1 dropdb simple_bank
migrateup:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v ./...
.PHONY: sqlc migrateup migratedown createdb postgres dropdb