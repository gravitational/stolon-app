FROM quay.io/gravitational/debian-grande:0.0.1

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN sed '/deb http:\/\/httpredir.debian.org\/debian\/ jessie-updates main contrib non-free/d' -i /etc/apt/sources.list && \
    echo "deb http://archive.debian.org/debian jessie-backports main contrib non-free" > /etc/apt/sources.list.d/jessie-backports.list && \
    echo "Acquire::Check-Valid-Until "false";" >> /etc/apt/apt.conf && \
    apt-get update && \
    apt-get install --yes --no-install-recommends -t jessie-backports \
        postgresql-9.6 \
        postgresql-client-9.6 \
        postgresql-contrib-9.6 && \
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

RUN chmod a+rx /usr/local/bin/run.sh
RUN useradd -ms /bin/bash stolon

EXPOSE 5431 5432 6431

RUN test -f /cleanup.sh && sh /cleanup.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
