# Stage 1: Build the statically linked Go binary
FROM golang:1.26-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /app

# Download dependencies first to cache this layer
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build a statically linked binary without debug info for minimum size
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o exporter ./cmd/exporter

# Stage 2: Minimal runtime image (distroless)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/exporter /exporter

# Run as non-root user
USER 65532:65532

ENTRYPOINT ["/exporter"]
