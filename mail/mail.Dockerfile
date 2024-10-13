FROM golang:1.22.2-alpine as builder
LABEL authors="arshia"


RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp


FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailApp /app

CMD ["/app/mailApp"]
