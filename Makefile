all: build

build:
	go build -o ./bin/server ./server/server.go
	go build -o ./bin/client ./client/client.go
	chmod +x ./bin/server
	chmod +x ./bin/client
