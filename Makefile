.PHONY: all
all: bin/nman bin/nman-daemon

bin/nman: **/**/*.go **/*.go
	@go build -o ./bin/nman cmd/nman/main.go

bin/nman-daemon: **/**/*.go **/*.go
	@go build -o ./bin/nman-daemon cmd/daemon/main.go
