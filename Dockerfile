FROM golang:1.16.10 as builder

WORKDIR /goplay

COPY pkg pkg
COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN apt update -y && apt install libasound2-dev -y
ENV GOPROXY=https://goproxy.io
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.14.3

COPY --from=builder /goplay/goplay /usr/local/bin/goplay

CMD ["/usr/local/bin/goplay"]
