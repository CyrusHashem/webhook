FROM alpine

RUN apk update && apk add ca-certificates

COPY webhook /usr/local/bin/webhook

USER 1

EXPOSE 8080

CMD webhook -listen :8080
