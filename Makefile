build:
	@go build -o tmp/rivate_theatre ./cmd/api/api.go

run: 
	@go run ./cmd/api/api.go