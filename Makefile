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

update-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init
dev:
	npx nodemon --exec go run main.go --signal SIGTERM