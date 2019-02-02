FROM quay.io/gravitational/debian-tall:stretch

ADD bin/stolonboot \
        /usr/local/bin/

ADD bootstrap.sh /bootstrap.sh
RUN chmod a+x /bootstrap.sh

CMD ["/bootstrap.sh"]
