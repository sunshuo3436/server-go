## Build tools developed with golang language

server:
	go build -o ./bin/server ./server/main.go

client:
	go build -o ./bin/client ./client/main.go

all: server client

.PHONY: server client