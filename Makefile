help:
	@echo "TODO"
	@echo "make start"

start:
	go run main.go --port 8080

deps:
	GO111MODULE=on go mod vendor

build:
	go build -o plex-watcher main.go


