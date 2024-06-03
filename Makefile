include .env

start:	
	go run ./cmd/api

migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext=sql -dir=migrations ${name}

migrations/up:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} up

migrations/down:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} down