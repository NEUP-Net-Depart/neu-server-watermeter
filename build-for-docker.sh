#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOFLAGS="" go build -a -installsuffix cgo -o main .