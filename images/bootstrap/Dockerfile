FROM quay.io/gravitational/debian-tall:0.0.1

ADD bin/stolonboot \
        /usr/local/bin/

ADD bootstrap.sh /bootstrap.sh
RUN chmod a+x /bootstrap.sh

CMD ["/bootstrap.sh"]

