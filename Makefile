postgres:
	docker run --name postgres-test-1 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest
createdb:
	docker exec -it postgres-test-1 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-test-1 dropdb simple_bank
migrateup:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/sherifzaher/clone-simplebank/db/sqlc Store
test:
	go test -v ./...
server:
	go run main.go
build:
	go build -o server main.go
run-prod:
	./server
proto:
	rm -f doc/swagger/*.swagger.json
	rm -f pb/*.go

	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge,merge_file_name=simple_bank \
    	proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans --host localhost --port 9090 -r repl

.PHONY: sqlc migrateup migratedown createdb postgres dropdb server mock proto evans