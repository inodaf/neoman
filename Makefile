.PHONY: all
all: ./bin/nman ./bin/nmand

./bin/nman: internal/**/*.go packages/**/*.go
	@go build -o ./bin/nman cmd/nman/main.go

./bin/nmand: internal/**/*.go packages/**/*.go
	@go build -o ./bin/nmand cmd/daemon/main.go
