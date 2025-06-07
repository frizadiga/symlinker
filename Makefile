.PHONY: install update tidy dev run clean build

BINARY_NAME=symlinker

install:
	go mod download

update:
	go get -u

tidy:
	go mod tidy

dev:
	go run -v .

start:
	./$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

build:
	go build -o $(BINARY_NAME) .

