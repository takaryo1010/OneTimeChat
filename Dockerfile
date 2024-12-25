# Goのインストール
FROM golang:1.23-alpine as go-builder

# Node.jsとGitのインストール
FROM node:23

# Gitをインストール
RUN apt-get update && apt-get install -y git

# Go環境を持ち込む
COPY --from=go-builder /usr/local/go /usr/local/go

# ポートの設定
EXPOSE 3000 8080


