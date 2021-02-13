FROM golang:1.15-alpine3.12

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN make build

CMD ["/app/bin/hello"]