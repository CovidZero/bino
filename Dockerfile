FROM golang:alpine3.11 as builder

COPY go.mod go.sum /usr/src/covid0-backend/
WORKDIR /usr/src/covid0-backend
RUN go mod download

COPY . /usr/src/covid0-backend
RUN go build -o /usr/local/bin/covid-api ./cmd/covid0-api

FROM alpine:3.11
COPY --from=builder /usr/local/bin/covid-api /usr/local/bin
RUN mkdir -p /var/covid0-api/temp-storage/covid0.db

EXPOSE 8080

CMD [ "/usr/local/bin/covid-api" ]
