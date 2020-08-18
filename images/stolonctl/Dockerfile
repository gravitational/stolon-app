FROM quay.io/gravitational/debian-grande:stretch

ARG CHANGESET
ARG VERSION

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm \
    RIG_CHANGESET=$CHANGESET \
    APP_VERSION=$VERSION


RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
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

ADD bin/stolonctl entrypoint.sh fix-permissions.sh pre-update.sh /usr/bin/
RUN chmod +x /usr/bin/stolonctl && \
    chmod +x /usr/bin/entrypoint.sh && \
    chmod +x /usr/bin/fix-permissions.sh && \
    chmod +x /usr/bin/pre-update.sh && \
    useradd -ms /bin/bash stolon

USER stolon
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
