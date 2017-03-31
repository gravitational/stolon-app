FROM quay.io/gravitational/debian-grande:0.0.1

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm \
    PGWEB_VERSION=0.9.6

RUN apt-get update && \
    apt-get install --yes --no-install-recommends curl unzip && \
    curl -sSL https://github.com/sosedoff/pgweb/releases/download/v${PGWEB_VERSION}/pgweb_linux_amd64.zip -o /pgweb_linux_amd64.zip && \
    unzip /pgweb_linux_amd64.zip -d /usr/local/bin/ && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/* \
        /pgweb_linux_amd64*

ADD rootfs/ /

RUN chmod a+rx /usr/local/bin/run.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
