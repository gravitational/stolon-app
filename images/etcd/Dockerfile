FROM quay.io/gravitational/debian-grande:buster as downloader
ARG ETCD_VERSION=v3.3.11
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -qq -y curl && \
    mkdir -p /etcd && \
    curl -L https://storage.googleapis.com/etcd/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz -o /etcd.tar.gz && \
    tar xzf /etcd.tar.gz -C /etcd --strip-components=1

FROM quay.io/gravitational/debian-tall:buster
COPY --from=downloader /etcd/etcd* /usr/bin/
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 2379

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/entrypoint.sh"]
