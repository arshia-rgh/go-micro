FROM alpine:latest
LABEL authors="arshia"

RUN mkdir /app

COPY mailApp /app

CMD ["/app/mailApp"]
