FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o app ./cmd/main.go
# 下载时区文件
RUN apk add tzdata


FROM scratch

COPY --from=builder /build/app /
COPY --from=builder /build/configs /
COPY --from=builder /build/certs /etc/ssl/certs/

# 拷贝时区文件
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 6565
CMD ["/app"]
