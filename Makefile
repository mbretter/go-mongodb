.PHONY: coverage test

test:
	go test -v ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

mocks:
	mockery