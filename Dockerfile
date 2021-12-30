# syntax=docker/dockerfile:1

# ---- Base ----
FROM golang:1.17 as base
RUN apt-get update
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... 

# ---- Dev ----
FROM base AS dev
CMD go run ./...

# ---- Test ----
FROM base AS test
CMD go test -count=1 -v ./...

# ---- Build ----
FROM base AS build
RUN go build -o go-api cmd/go-api/main.go

# ---- Release ----
FROM debian:buster-slim AS release
COPY --from=build /go/src/app/go-api ./
CMD ["./go-api"]