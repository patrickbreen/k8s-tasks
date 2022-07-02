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
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /canary /canary
EXPOSE 9000
# Run the hello binary.
ENTRYPOINT ["/app"]