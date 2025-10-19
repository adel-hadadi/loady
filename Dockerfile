FROM golang:1.23 AS builder

WORKDIR /src

COPY go.* .

RUN go mod download

COPY . .

ENV DOCKER_API_VERSION=1.39

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loady ./cmd/main.go

FROM alpine:latest

COPY --from=builder /src/loady loady

RUN mkdir -p /etc/loady

CMD [ "./loady" ]