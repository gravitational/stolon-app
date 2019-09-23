FROM quay.io/gravitational/debian-grande:stretch

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        gettext \
        hostname \
        procps \
        pgbouncer \
        postgresql-common \
        postgresql-client-9.6 \
        openssl && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/*

RUN useradd -ms /bin/bash stolon && \
    chmod +x /usr/local/bin/init-container.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

ADD init-container.sh /usr/local/bin/init-container.sh
ADD entrypoint.sh /usr/local/bin/entrypoint.sh
ADD common_lib.sh /common_lib.sh
ADD pgbouncer_auth.sql /home/stolon/pgbouncer_auth.sql

USER stolon
EXPOSE 6432
ENTRYPOINT ["dumb-init", "--", "/usr/local/bin/entrypoint.sh"]
