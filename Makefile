.PHONY: swagger test

swagger:
	swag init -g cmd/app/main.go -o docs

test:
	go test -v ./internal/service/...
