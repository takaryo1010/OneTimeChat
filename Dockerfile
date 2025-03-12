# Ubuntu 22.04 (LTS)をベースに使用
FROM ubuntu:22.04

# ミラーサーバーを変更し、より高速なサーバーを使用
RUN sed -i 's/archive.ubuntu.com/ftp.jaist.ac.jp/g' /etc/apt/sources.list

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    language-pack-ja \
    locales \
    && apt-get clean

# bashの日本語化
RUN locale-gen

RUN export LC_ALL=ja_JP.utf8

ENV LANG=ja_JP.UTF-8

# NodeSourceのリポジトリを追加して最新のNode.jsをインストール
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs

# Go環境をインストール
RUN curl -LO https://golang.org/dl/go1.23.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xvzf go1.23.0.linux-amd64.tar.gz \
    && rm go1.23.0.linux-amd64.tar.gz

# Goの環境変数を設定
ENV PATH="/usr/local/go/bin:${PATH}"

# GOPATHの設定
ENV GOPATH="/root/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Goのリンターをインストール
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${GOPATH}/bin v1.64.7
RUN go install golang.org/x/tools/cmd/goimports@latest
# 作業ディレクトリを作成
WORKDIR /workspace

# リポジトリをクローン
RUN git clone https://github.com/takaryo1010/OneTimeChat.git .

# ポートの設定
EXPOSE 3000 8080
