@echo off
set CGO_ENABLED=0
set GOOS=linux
set GOFLAGS=
go build -a -installsuffix cgo -o main .
