 # Set the Go version as a build argument
ARG GO_VERSION=1.23

# Stage 1: Build the Go application
FROM golang:${GO_VERSION}-alpine AS builder
ARG SERVICE_PORT="80"
ARG TARGETARCH=amd64

RUN apk add --no-cache tzdata

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=$TARGETARCH
WORKDIR /app

# Copy the Go application source code
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -installsuffix cgo -ldflags '-s -w' -o main ./main.go


FROM alpine:3.17
ARG SERVICE_PORT="80"
ENV PORT=${SERVICE_PORT}

# Install necessary packages
RUN apk add --no-cache \
    bash \
    gettext \
    curl \
    tzdata \
    openssl \
    postgresql-client


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
COPY --from=builder /app/script/. ./script/.
COPY --from=builder /app/src/. ./src/.
# Expose the specified service port
EXPOSE $PORT

# Run the compiled Go binary
CMD ["./main"]
