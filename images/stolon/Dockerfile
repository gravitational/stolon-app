FROM quay.io/gravitational/debian-grande:stretch

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        postgresql-9.6 \
        postgresql-client-9.6 \
        postgresql-contrib-9.6

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
