FROM quay.io/gravitational/debian-grande:0.0.1

ENV STSENTINEL_INITIAL_CLUSTER_CONFIG=/etc/stolon/cluster-config.json

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
	postgresql \
	postgresql-client \
    postgresql-contrib && \
    apt-get clean

ADD run.sh bin/stolon-keeper bin/stolon-sentinel bin/stolon-proxy /usr/local/bin/

ADD cluster-config.json /etc/stolon/cluster-config.json

RUN chmod +x /usr/local/bin/stolon-keeper /usr/local/bin/stolon-sentinel /usr/local/bin/stolon-proxy /usr/local/bin/run.sh

RUN useradd -ms /bin/bash stolon

EXPOSE 5431 5432 6431

ENTRYPOINT ["/usr/local/bin/run.sh"]
