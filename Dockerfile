# Use buildx for multi-platform builds
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/lynk-mcp"

RUN apk add --no-cache make git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build for multiple architectures
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o lynk-mcp ./cmd/lynk-mcp

RUN chmod +x lynk-mcp

# Final stage
FROM alpine:3.21
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/lynk-mcp"
LABEL org.opencontainers.image.description="MCP server for Lynk version management API"
LABEL org.opencontainers.image.licenses=Apache-2.0

COPY --from=builder /app/lynk-mcp /app/lynk-mcp

# Disable version check
ENV INTERLYNK_DISABLE_VERSION_CHECK=true

ENTRYPOINT ["/app/lynk-mcp"]
