FROM quay.io/gravitational/debian-grande:buster as downloader

ARG GOLANGCI_LINT_VER

RUN apt-get update && apt-get install wget -yy && \
    wget https://github.com/golangci/golangci-lint/releases/download/v$GOLANGCI_LINT_VER/golangci-lint-$GOLANGCI_LINT_VER-linux-amd64.tar.gz && \
	tar -xvf golangci-lint-$GOLANGCI_LINT_VER-linux-amd64.tar.gz \
    golangci-lint-$GOLANGCI_LINT_VER-linux-amd64/golangci-lint --strip-components=1

FROM golang:1.13-buster
COPY --from=downloader /golangci-lint /usr/local/bin/
