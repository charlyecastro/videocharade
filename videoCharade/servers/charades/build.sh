#! /bin/bash
GOOS=linux go build
docker build -t charlyecastro/charades .
go clean
