FROM golang:1.22.3-alpine AS build

RUN apk update && \
    apk add --update gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./

# Building this takes a long time - so build it in an earlier layer
RUN go install github.com/mattn/go-sqlite3

COPY . .

RUN CGO_ENABLED=1 go build -ldflags="-s -w"

FROM alpine

WORKDIR /

COPY --from=build /build/master-farmer .

ENTRYPOINT ["./master-farmer"]