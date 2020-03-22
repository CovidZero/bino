.PHONY: build tidy test package watch devsetup test run precommit setupprecommit setupmodd testscompile

build:
	go build ./...

# apenas checando se os tests compilam, pois como rodamos todos os testes o tempo todo
# esperar todo o ambiente de teste ser carregado para depois receber erros de compilação
testscompile:
	go test -run='^$$' ./...

test: build testscompile
	bash scripts/test/integration_test/test.sh

tidy: test
	go mod tidy
	go fmt ./...

watch: build
	modd

package: tidy
	docker build . -t covidzero/bino:latest

run: package
	docker run -p '8080:8080' -e COVID0_TEMP_BUCKET -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY covidzero/bino:latest

precommit: tidy

devsetup: setupmodd setupprecommit

setupmodd:
	GO111MODULE=on go get github.com/cortesi/modd/cmd/modd
	go install github.com/cortesi/modd/cmd/modd

setupprecommit:
	chmod u+x scripts/githooks/pre-commit.sh
	cp scripts/githooks/pre-commit.sh .git/hooks/pre-commit
