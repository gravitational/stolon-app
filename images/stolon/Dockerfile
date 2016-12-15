FROM quay.io/gravitational/debian-grande:0.0.1

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        postgresql \
        postgresql-client \
        postgresql-contrib && \
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
