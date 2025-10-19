FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.* .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:latest

COPY --from=builder /src/app app

CMD [ "./app" ]
