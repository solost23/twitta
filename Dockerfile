FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o app ./cmd/twitta/main.go
# 下载时区文件
RUN apk add tzdata


FROM scratch

COPY --from=builder /build/app /
#COPY --from=builder /build/configs /

# 拷贝时区文件
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# 设置时区
ENV TZ=Asia/Shanghai
# 安装证书，防止Unknown desc = x509: certificate signed by unknown authority错误
# 注意： 安装服务器的证书
ADD ca-certificates.crt /etc/ssl/certs/

EXPOSE 6565
CMD ["/app"]
