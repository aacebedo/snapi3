# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.140.1/containers/alpine/.devcontainer/base.Dockerfile

ARG VARIANT="20.04"
FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu-${VARIANT}

# ** [Optional] Uncomment this section to install additional packages. **
RUN apt update \
     && apt install golang zsh libgtk-3-dev libcairo2-dev libglib2.0-dev -y
WORKDIR /usr/local
RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

USER vscode
ENV PATH="/home/vscode/go/bin:${PATH}"
ENV GOPATH="/home/vscode/go"
RUN GO111MODULE=on go get -v -u github.com/go-critic/go-critic/cmd/gocritic
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0
RUN GO111MODULE=on go get mvdan.cc/gofumpt
RUN go get github.com/daixiang0/gci
RUN go get golang.org/x/tools/cmd/goimports

