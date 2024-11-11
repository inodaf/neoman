.PHONY: all
all: bin/nman bin/nmand

bin/nman: **/**/*.go **/*.go
	@go build -o ./bin/nman cmd/nman/main.go

bin/nmand: **/**/*.go **/*.go
	@go build -o ./bin/nmand cmd/daemon/main.go
