# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:latest AS build

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./cmd

##
## Deploy
##
FROM scratch

WORKDIR /

COPY --from=build /service /service

#RUN apk add --no-cache \
#        musl
#
ENTRYPOINT ["/service"]

