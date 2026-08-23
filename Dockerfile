# Stage 1: Build the statically linked Go binary
FROM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder

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
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7

WORKDIR /

COPY --from=builder /app/exporter /exporter

# Run as non-root user
USER 65532:65532

ENTRYPOINT ["/exporter"]
