FROM quay.io/gravitational/debian-grande:0.0.1

ARG DEBIAN_FRONTEND=noninteractive

ENV TERM=xterm \
    PGBINOLD=/usr/lib/postgresql/9.4/bin \
    PGBINNEW=/usr/lib/postgresql/9.6/bin \
    PG_VERSION_OLD=9.4 \
    PG_VERSION_NEW=9.6

RUN echo "deb http://deb.debian.org/debian jessie-backports main contrib non-free" > /etc/apt/sources.list.d/jessie-backports.list && \
    apt-get update && \
    apt-get install --yes --no-install-recommends \
        curl \
        rsync \
        postgresql-9.4 \
        postgresql-contrib-9.4 && \
    apt-get install --yes --no-install-recommends -t jessie-backports \
        postgresql-9.6 \
        postgresql-client-9.6 \
        postgresql-contrib-9.6 \
        s3cmd && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/*

RUN useradd -ms /bin/bash stolon
ADD keeper-upgrade.sh /usr/local/bin/
ADD clean-postgres-data.sh /usr/local/bin/
ADD move-postgres-data.sh /usr/local/bin/
ADD entrypoint.sh /
RUN chmod +x /usr/local/bin/keeper-upgrade.sh && \
    chmod +x /usr/local/bin/clean-postgres-data.sh && \
    chmod +x /usr/local/bin/move-postgres-data.sh

USER stolon
ENTRYPOINT ["dumb-init", "/entrypoint.sh"]
