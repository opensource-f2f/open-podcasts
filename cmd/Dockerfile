# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api api/
COPY pkg pkg/
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 go build -a -o yaml-rss cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source=https://github.com/opensource-f2f/open-podcasts-yaml-rss
WORKDIR /
COPY --from=builder /workspace/yaml-rss /usr/bin/yaml-rss
USER 65532:65532

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
