#! /bin/bash
GOOS=linux go build
docker build -t charlyecastro/pagesummaryapi .
go clean