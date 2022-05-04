.PHONY:all
all:build

.PHONY:build
build:
	go build

.PHONY:srv
srv:
	go-tour.exe server


.PHONY:cli
cli:
	go-tour.exe client

.PHONY: test
test:
	go-tour.exe server --gateway="127.0.0.1:50051"