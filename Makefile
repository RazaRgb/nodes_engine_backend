.DEFAULT_GOAL := build

.PHONY:fmt vet build clean run
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build: vet
	go build -o backend ./src
clean:
	go clean
	rm -rf ./backend
run: clean build
	./backend
