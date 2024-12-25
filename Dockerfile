# Node.js用のベースイメージを使用
FROM node:23

# Gitをインストール
RUN apt-get update && apt-get install -y git

# Go環境をインストール
RUN curl -LO https://golang.org/dl/go1.23.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xvzf go1.23.0.linux-amd64.tar.gz \
    && rm go1.23.0.linux-amd64.tar.gz

# Goの環境変数を設定
ENV PATH="/usr/local/go/bin:${PATH}"

# ポートの設定
EXPOSE 3000 8080

