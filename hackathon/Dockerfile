# ビルド用イメージ
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server

# 実行用イメージ
FROM alpine:latest
# Asia/Tokyo を使うために tzdata をインストール
RUN apk add --no-cache tzdata
# タイムゾーンを Asia/Tokyo に設定（環境変数とリンク）
ENV TZ=Asia/Tokyo
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
