FROM golang:1.19-alpine AS build-stage

RUN apk add upx
WORKDIR /messages-src
COPY . .
RUN go build -o /messages .
RUN upx /messages

FROM alpine:latest AS release-stage

COPY --from=build-stage /messages /messages
# -- Environment variables
ENV MEMPHIS_STATION     "messages"
ENV MEMPHIS_CONSUMER    "messages"
ENV MEMPHIS_HOST        "memphis"
ENV MEMPHIS_USERNAME    "root"
ENV MEMPHIS_PASSWORD    "memphis"
ENV MEMPHIS_CONN_TOKEN  ""
ENV POSTGRES_DSN        "host=postgres user=sulcud password=sulcud dbname=sulcud port=5432 sslmode=disable"
ENV REDIS_ADDR          "redis:6379"
ENV REDIS_DB            "1"
ENV SMTP_FROM           "messages@mal.com"
ENV SMTP_HOST           "mailhog"
ENV SMTP_PORT           "1025"
ENV SMTP_USERNAME       ""
ENV SMTP_PASSWORD       ""
ENV SMTP_NO_TLS         "1"
# -- Environment variables
ENTRYPOINT [ "sh", "-c", "/messages" ]