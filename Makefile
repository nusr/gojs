build:
	go build
lint:
	go vet ./...
  go fmt ./...
  go mod tidy
test:
	go test ./...
coverage:
	go test ./... -coverprofile=coverageFile
	go tool cover -html=coverageFile