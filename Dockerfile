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
    openssh-server \
    && apt-get clean

# bashの日本語化
RUN locale-gen
RUN export LC_ALL=ja_JP.utf8
ENV LANG=ja_JP.UTF-8

# SSHの設定
RUN mkdir /var/run/sshd && \
    echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config && \
    echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config

# 環境変数から root のパスワードを設定
ARG ROOT_PASSWORD
RUN echo "root:${ROOT_PASSWORD}" | chpasswd

# NodeSourceのリポジトリを追加して最新のNode.jsをインストール
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs

# Go環境をインストール
RUN curl -LO https://golang.org/dl/go1.23.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xvzf go1.23.0.linux-amd64.tar.gz \
    && rm go1.23.0.linux-amd64.tar.gz

# Goの環境変数を.bashrcに追加
RUN echo 'export PATH="/usr/local/go/bin:$PATH"' >> /root/.bashrc
RUN echo 'export GOPATH="/root/go"' >> /root/.bashrc
RUN echo 'export PATH="${GOPATH}/bin:$PATH"' >> /root/.bashrc


 # Goのリンターをインストール
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${GOPATH}/bin v1.64.7
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN echo 'gofmt -w . && goimports -w . && golangci-lint run' >> /root/.bashrc
# 作業ディレクトリを作成
WORKDIR /workspace

# SSH ログイン時に /workspace に移動するように設定
RUN echo 'cd /workspace' >> /root/.bashrc

# リポジトリをクローン
RUN git clone https://github.com/takaryo1010/OneTimeChat.git .

# ポートの設定
EXPOSE 22 3000 8080

# SSH デーモンを起動
CMD ["/bin/bash", "-c", "source /root/.bashrc && /usr/sbin/sshd -D"]
