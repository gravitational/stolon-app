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

ADD bin/ /usr/local/bin/
ADD rootfs/ /

RUN chmod a+rx /usr/local/bin/run.sh && \
    chmod a+rx /usr/local/bin/cron-wal-e.sh

RUN useradd -ms /bin/bash stolon
RUN mkfifo --mode 0666 /var/log/cron.log

EXPOSE 5431 5432 6431

RUN test -f /cleanup.sh && sh /cleanup.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
