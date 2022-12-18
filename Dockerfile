FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod tidy
RUN go build -o read-only-service ./cmd/main.go

CMD ["./read-only-service"]
