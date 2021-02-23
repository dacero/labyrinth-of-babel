FROM golang:1.15-alpine3.12

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

RUN make build

CMD ["/app/bin/lob"]