#!/usr/bin/env bash
set GOOS=linux
set GOARCH=amd64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -v" -o main main.go
chmod +777 main
sam local start-api