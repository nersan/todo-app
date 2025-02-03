# ベースイメージ
FROM golang:1.20

# 作業ディレクトリの設定
WORKDIR /app

# モジュールファイルのコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN go build -o main ./cmd/main.go

# アプリケーションを実行
CMD ["/app/main"]
