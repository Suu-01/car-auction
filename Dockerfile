# 1) ビルドステージ
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 必要なシステムパッケージをインストール (CA 証明書, git など)
RUN apk add --no-cache ca-certificates git

# モジュールファイルのみコピーして依存関係をダウンロード (キャッシュ活用)
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションソースをコピー
COPY . .

# アプリケーションをビルド (CGO を無効化)
RUN CGO_ENABLED=0 \
    go build -o car-auction ./cmd/server

# 2) 実行ステージ (最小イメージ)
FROM scratch

# ルート証明書をコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ビルド済みバイナリをコピー
COPY --from=builder /app/car-auction /car-auction

# デフォルトでリッスンするポート
EXPOSE 8080

# コンテナ起動時に実行されるコマンド
ENTRYPOINT ["/car-auction"]
