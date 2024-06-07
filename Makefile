include .env

build:	
	go build -o ./build.exe ./cmd/api

run:
	./build.exe

build/run:
	make build && make run

migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext=sql -dir=migrations ${name}

migrations/up:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} up

migrations/down:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} down