bin/nman: cmd/**/*.go internal/**/*.go
	@go build -o ./bin/nman cmd/nman/main.go
