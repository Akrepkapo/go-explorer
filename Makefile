export GOPRIVATE=lib.venas.io
export GO111MODULE=on

build:
	go get -u
	go mod tidy -v
	go build

initdb:
	go-ibax-explorer initDatabase
start:
	go-ibax-explorer start

all: build

startup: initdb start
