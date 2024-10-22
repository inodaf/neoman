bin/nman: **/**/*.go **/*.go
	@go build -o ./bin/nman cmd/nman/main.go
