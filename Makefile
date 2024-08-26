.PHONY: migrate
migrate:
	goose -dir=migrations sqlite3 eventor.db up

.PHONY: run
run:
	go run cmd/eventor/main.go
