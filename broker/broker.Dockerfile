FROM alpine:latest
LABEL authors="arshia"

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]