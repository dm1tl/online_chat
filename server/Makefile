postgresinit:
	docker run --rm --name postgres_db -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres

postgres:
	docker exec -it postgres_db psql

createdb:
	docker exec -it postgres_db createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres_db dropdb go-chat

migrateup:
	migrate -path migrations/chat_schema -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations/chat_schema -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown