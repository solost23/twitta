FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cp configs/config.docker.yml configs/config.yml
RUN go build -o app ./cmd/main.go
RUN apk add tzdata


FROM scratch

COPY --from=builder /build/app /
COPY --from=builder /build/configs /configs
COPY --from=builder /build/certs /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

EXPOSE 6565
CMD ["/app", "-d", "/"]
