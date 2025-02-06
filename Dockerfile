 # Set the Go version as a build argument
ARG GO_VERSION=1.23

# Stage 1: Build the Go application
FROM golang:${GO_VERSION}-alpine AS builder
ARG SERVICE_PORT="80"
RUN apk add --no-cache tzdata

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=$TARGETARCH
WORKDIR /app

# Copy the Go application source code
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -installsuffix cgo -ldflags '-s -w' -o main ./main.go

# Stage 2: Use wkhtmltopdf as a base
FROM surnet/alpine-wkhtmltopdf:3.16.2-0.12.6-full AS wkhtmltopdf

# Stage 3: Create final image
FROM alpine:3.14
ARG SERVICE_PORT="80"
ENV PORT=${SERVICE_PORT}

# Install necessary packages
RUN apk add --no-cache \
    bash=5.1.16-r0 \
    gettext=0.21-r0 \
    curl=8.0.1-r0 \
    xvfb=1.20.11-r0 \
    tzdata=2023c-r0 \
    openssl=1.1.1t-r2 \
    postgresql-client=16.12-r0

# Set the working directory
WORKDIR /app

# Create an empty .env file
RUN touch .env

# Add User
RUN addgroup -S application
RUN adduser -S user -G application
RUN chown -R user:application /app
USER user

# Copy the compiled Go binary and other resources
COPY --from=builder /app/main ./main
COPY --from=builder /app/scripts/. ./scripts/.

# Expose the specified service port
EXPOSE $PORT

# Run the compiled Go binary
CMD ["./main"]
