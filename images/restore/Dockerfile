FROM quay.io/gravitational/debian-grande:0.0.1

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        daemontools file \
        gcc libc6-dev \
        python3-pip python3-dev \
        lzop pv cron gawk \
        postgresql \
        postgresql-client \
        postgresql-contrib && \
    pip3 install wal-e boto && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/*

ADD rootfs/ /

RUN chmod a+rx /usr/local/bin/run.sh

EXPOSE 5431 5432 6431

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
