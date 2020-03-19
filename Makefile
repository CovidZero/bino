.PHONY: build tidy test package watch devsetup

build:
	go build ./...


test: build
	go test ./...

tidy: build
	go mod tidy
	go fmt ./...

watch: build
	modd

package:
	docker build . -t andrebq/covid0-backend:latest

devsetup:
	mv go.mod go.mod
	GO111MODULE=on go get github.com/cortesi/modd/cmd/modd
	go install github.com/cortesi/modd/cmd/modd
