FROM alpine:latest
LABEL authors="arshia"


RUN mkdir /app

COPY authApp /app

CMD ["/app/authApp"]
