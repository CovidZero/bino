.PHONY: build tidy test package watch devsetup integration_test

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
	docker build . -t covidzero/bino:latest

run: package
	docker run -p '8080:8080' -e COVID0_TEMP_BUCKET -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY covidzero/bino:latest

devsetup:
	mv go.mod go.mod
	GO111MODULE=on go get github.com/cortesi/modd/cmd/modd
	go install github.com/cortesi/modd/cmd/modd

integration_test:
	cd test-compose && docker-compose up -d
	go run internal/testutil/waitfor/main.go -a "localhost:4572"
