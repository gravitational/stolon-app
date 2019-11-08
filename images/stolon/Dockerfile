ARG PGVERSION=9.6
# base build image
FROM golang:1.13-buster AS build_base
ARG STOLON_REVISION=656bc2fe5d801193b412bcfb57dbfa84cfd9f386

RUN set -eux; \
    apt-get update && \
    apt-get install -qq -y git && \
    git clone https://github.com/gravitational/stolon.git /stolon && \
    cd /stolon && git checkout $STOLON_REVISION

WORKDIR /stolon

RUN go mod download

#######
####### Build the stolon binaries
#######
FROM build_base AS builder

RUN make

#######
####### Download etcd binaries
#######
FROM quay.io/gravitational/debian-grande:buster as downloader
ARG ETCD_VERSION=v3.3.11
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -qq -y curl && \
    mkdir -p /etcd && \
    curl -L https://storage.googleapis.com/etcd/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz -o /etcd.tar.gz && \
    tar xzf /etcd.tar.gz -C /etcd --strip-components=1

#######
####### Build the final image
#######
FROM postgres:$PGVERSION

ADD rootfs/ /

RUN apt-get update && \
    apt-get install dumb-init && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/* && \
    chmod +x /usr/local/bin/init-container.sh && \
    useradd -ms /bin/bash stolon

EXPOSE 5432

COPY --from=builder /stolon/bin/ /usr/local/bin
COPY --from=downloader /etcd/etcd* /usr/bin/

RUN chmod +x /usr/local/bin/stolon-keeper /usr/local/bin/stolon-sentinel /usr/local/bin/stolon-proxy /usr/local/bin/stolonctl

USER stolon

ENTRYPOINT ["stolon-keeper"]
