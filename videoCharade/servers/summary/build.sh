#! /bin/bash
GOOS=linux go build
docker build -t charlyecastro/summary .
go clean