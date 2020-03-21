#!/usr/bin/env bash

# TODO: incluir um target no Makefile para validar esse script usando shellcheck (via docker)
# provavelmente deixei algo passar

set -euo pipefail

readonly initialFolder=$(pwd)
readonly composeFolder="${initialFolder}/test-compose/"

function docker-compose-down {
    cd "${composeFolder}"
    docker-compose down
}

function docker-compose-up {
    cd "${composeFolder}"
    docker-compose up -d
}

function go-test {
    cd "${initialFolder}"
    export COVID0_TEMP_BUCKET='s3://mybucket?region=some&endpoint=localhost:4572&disableSSL=true&s3ForcePathStyle=true'
    export AWS_ACCESS_KEY_ID='keyid'
    export AWS_SECRET_ACCESS_KEY='keyval'
    go test ./...
}

function wait-for-localstack {
    go run "${initialFolder}/scripts/test/waitfor/main.go" -a "localhost:4572"
}

trap 'docker-compose-down' EXIT

docker-compose-up
wait-for-localstack
go-test
docker-compose-down

trap '' EXIT
