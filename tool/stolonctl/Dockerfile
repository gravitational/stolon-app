FROM quay.io/gravitational/debian-grande:0.0.1

RUN apt-get update && apt-get install curl net-tools -y

ADD build/* /usr/bin/
ADD entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/stolonctl && \
    chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
