bin/nman: $(shell ls **/*.go)
	@go build -o ./bin/nman cmd/nman/main.go
