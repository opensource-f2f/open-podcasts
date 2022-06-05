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
COPY cmd cmd/
COPY api api/
COPY pkg pkg/
COPY entrypoint.sh entrypoint.sh
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 go build -a -o yaml-rss cmd/main.go

FROM alpine:3.10
LABEL org.opencontainers.image.source=https://github.com/opensource-f2f/open-podcasts
WORKDIR /
COPY --from=builder /workspace/yaml-rss /usr/bin/yaml-rss
COPY --from=builder /workspace/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
