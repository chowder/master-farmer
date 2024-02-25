FROM golang:1.22.0-alpine3.19 AS build

RUN apk update && \
    apk add --update gcc musl-dev

WORKDIR /build

COPY . .

RUN CGO_ENABLED=1 go build -ldflags="-s -w"

FROM alpine:3.19

WORKDIR /

COPY --from=build /build/master-farmer .

ENTRYPOINT ["./master-farmer"]
