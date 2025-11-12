# Simple Makefile for a Go project

default: run

# Run the application
run:
	go run ./main.go

build:
	echo "Building..."
	go build -o bin/main.exe main.go

install:
	echo "Installing..."
	go install