############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /build
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
WORKDIR /build/canary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /canary
############################
# STEP 2 build a small image
# I use alpine instead of scratch, because I need the basic OS TLS CA's
############################
FROM alpine:3.16.0
# Copy our static executable.
COPY --from=builder /canary /canary
COPY canary/certs/client.crt /client.crt
COPY canary/certs/client.key /client.key
EXPOSE 9000
# Run the hello binary.
ENTRYPOINT ["/canary"]
