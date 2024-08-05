
installSQLC:
	sudo snap install sqlc

runDB:
	docker run --name bank -p 5432:5432 -e POSTGRES_USER=bank -e POSTGRES_PASSWORD=bank -e POSTGRES_DB=bank -d postgres:16-alpine

startDB:
	docker start bank

stopDB:
	docker stop bank

resetDB: dropDB createDB migrateUp

dropDB:
	docker exec -it bank dropdb -U bank bank

createDB:
	docker exec -it bank createdb -U bank --owner bank bank

migrateUp:
	migrate -path db/migration -database "postgresql://bank:bank@localhost:5432/bank?sslmode=disable" --verbose up

migrateDown:
	migrate -path db/migration -database "postgresql://bank:bank@localhost:5432/bank?sslmode=disable" --verbose down 1

## Create new migration file with name `name` ex: make migrateNew name=create_table_user
migrateNew:
	migrate create -ext sql -dir db/migration -seq $(name)

## Generate the go files from the sqlc queries
sqlc:
	sqlc generate

test:
	GOFLAGS="-count=1" go test -v -cover ./...

server:
	go run main.go

## Generate the mock files for the db
mock:
	mockgen -destination db/mock/store.go -package mockdb  github.com/RenanWinter/bank/db/sqlc Store

.PHONY: runDB startDB stopDB resetDB dropDB createDB migrateUP migrateDown migrateNew sqlc test server mock