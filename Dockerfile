# ---- Build Stage ----
# Use the official Golang image as the base image for building the application.
# Specify the version and OS variant requested.
FROM golang:1.24.1-alpine3.21 AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files first. This leverages Docker's layer caching.
# If go.mod and go.sum haven't changed, Docker will reuse the dependency layer.
COPY go.mod go.sum ./

# Download the Go module dependencies.
# Using "go mod download" explicitly downloads dependencies without trying to build.
RUN go mod download

# Copy the rest of the application source code into the container.
COPY . .

# Build the Go application.
# -o /app/binary specifies the output filename and path for the compiled binary.
# CGO_ENABLED=0 builds a statically linked binary (no external C dependencies),
# which is often preferred for minimal container images.
# -ldflags="-w -s" strips debugging information, reducing the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/binary .

# ---- Final Stage ----
FROM alpine:3.21.3

WORKDIR /app

# Copy the statically built binary from the builder stage
COPY --from=builder /app/binary /app/binary

# Run the binary as the container's entrypoint
ENTRYPOINT ["/app/binary"]
