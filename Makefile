build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/image-srv ./cmd/image-srv/
test:
	go test ./...