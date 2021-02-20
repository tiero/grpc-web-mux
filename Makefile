.PHONY: fmt server

## fmt: Go Format
fmt:
	@echo "Gofmt..."
	@if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted"; exit 1; fi

## server: Runs an insecure server as an example
server: 
	go run ./cmd/server/* -insecure