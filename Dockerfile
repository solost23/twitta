FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \

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

EXPOSE 6565
CMD ["/app"]