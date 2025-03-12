# Ubuntu 22.04 (LTS)をベースに使用
FROM ubuntu:22.04

# ミラーサーバーを変更し、より高速なサーバーを使用
RUN sed -i 's/archive.ubuntu.com/ftp.jaist.ac.jp/g' /etc/apt/sources.list

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    nodejs \
    npm \
    && apt-get clean

# Go環境をインストール
RUN curl -LO https://golang.org/dl/go1.23.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xvzf go1.23.0.linux-amd64.tar.gz \
    && rm go1.23.0.linux-amd64.tar.gz

# Goの環境変数を設定
ENV PATH="/usr/local/go/bin:${PATH}"

# Goのリンターをインストール
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.7


# Goのパスを設定
RUN echo 'export PATH=$(go env GOPATH)/bin:$PATH' >> ~/.bashrc

# 作業ディレクトリを作成
RUN mkdir ./workspace

# ポートの設定
EXPOSE 3000 8080

# 作業ディレクトリに移動
WORKDIR /workspace

