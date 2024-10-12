FROM golang:1.22.2-alpine as builder
LABEL authors="arshia"

RUN mkdir /app

COPY . /app

RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

RUN chmod +x /app/loggerApp


FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/loggerApp /app

CMD ["/app/loggerApp"]