FROM quay.io/gravitational/debian-grande:0.0.1

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    echo "deb http://deb.debian.org/debian jessie-backports main contrib non-free" > /etc/apt/sources.list.d/jessie-backports.list && \
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

RUN chmod a+rx /usr/local/bin/run.sh && \
    chmod a+rx /usr/local/bin/fix-certificates.sh
RUN useradd -ms /bin/bash stolon && \
    mkdir -p /home/stolon/bin && \
    mkdir -p /home/stolon/secrets && \
    chown stolon:stolon -R /home/stolon

EXPOSE 5431 5432 6431

RUN test -f /cleanup.sh && sh /cleanup.sh

USER stolon
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
