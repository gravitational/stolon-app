FROM quay.io/gravitational/debian-grande:stretch

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        hostname \
        procps \
        openssl && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        ~/.bashrc \
        /usr/share/doc/ \
        /usr/share/doc-base/ \
        /usr/share/man/ \
        /tmp/*

ADD init-container.sh /usr/local/bin/init-container.sh
RUN useradd -ms /bin/bash -u 65533 stolon && \
    chmod +x /usr/local/bin/init-container.sh

ENTRYPOINT ["dumb-init", "--", "/usr/local/bin/init-container.sh"]
