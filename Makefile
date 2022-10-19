build: 
	go build
test: 
	go test ./...
test-js:
	go build && ./gojs testData.js
lint:
	go vet ./... && go fmt ./... && go mod tidy
coverage: 
	go test ./... -coverprofile=coverageFile && go tool cover -html=coverageFile