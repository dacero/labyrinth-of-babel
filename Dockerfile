FROM golang:1.15-alpine3.12 as build
RUN apk add --no-cache build-base bash
WORKDIR /app
COPY . .
RUN make build
CMD ["/app/bin/lob"]