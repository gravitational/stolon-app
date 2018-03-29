FROM quay.io/gravitational/debian-grande:0.0.1

ADD bin/stolonctl /usr/bin/stolonctl
ADD entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/stolonctl && \
    chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
