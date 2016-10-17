FROM quay.io/gravitational/debian-tall:0.0.1

ADD bin/cfssl \
    bin/cfssljson \
    bin/stolonboot \
        /usr/local/bin/
RUN chmod a+x /usr/local/bin/cfssl*

ADD bootstrap.sh /bootstrap.sh
RUN chmod a+x /bootstrap.sh
CMD ["/bootstrap.sh"]

