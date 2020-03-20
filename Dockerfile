FROM golang:alpine3.11 as builder

COPY go.mod go.sum /usr/src/bino/
WORKDIR /usr/src/bino
RUN go mod download

COPY . /usr/src/bino
RUN go build -o /usr/local/bin/bino ./cmd/bino

FROM alpine:3.11
COPY --from=builder /usr/local/bin/bino /usr/local/bin

EXPOSE 8080

CMD [ "/usr/local/bin/bino" ]
