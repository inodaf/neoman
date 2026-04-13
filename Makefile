.PHONY: all
all: ./bin/nman ./bin/nmand

./bin/nman: ./internal/app/**/*.go
	@go build -o ./bin/nman ./cmd/nman

./bin/nmand: ./internal/daemon/**/*.go
	@go build -o ./bin/nmand ./cmd/daemon

.PHONY: migration
migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration NAME=description"; \
		echo "Example: make migration NAME=add_user_table"; \
		exit 1; \
	fi
	@go run github.com/pressly/goose/v3/cmd/goose@latest -dir cmd/daemon/db/migrations sqlite3 ./dev.db create $(NAME) sql
	@echo "Migration created: cmd/daemon/db/migrations/*_$(NAME).sql"
