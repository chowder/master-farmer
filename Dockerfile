FROM golang:1.22.0-alpine3.19 AS build

WORKDIR /build

COPY . .

RUN go build

FROM alpine:3.19

WORKDIR /

COPY --from=build /build/master-farmer .

ENTRYPOINT ["./master-farmer"]
