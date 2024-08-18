clean:
	rm -rf bin/*

build:
	go build -o bin/webapp main.go

image:
	docker build -t ghcr.io/bryopsida/gofiber-pug-starter:local .

test:
	go test -v ./...
	
lint:
	go install golang.org/x/lint/golint@latest
	golint ./...
	go vet ./...