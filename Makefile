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