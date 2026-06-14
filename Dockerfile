# Stage 1: Build the statically linked Go binary
FROM golang:1.26-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c AS builder

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
FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240

WORKDIR /

COPY --from=builder /app/exporter /exporter

# Run as non-root user
USER 65532:65532

ENTRYPOINT ["/exporter"]
