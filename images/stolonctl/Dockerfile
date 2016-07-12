FROM quay.io/gravitational/debian-tall:0.0.1

ADD bin/stolonctl /usr/local/bin/stolonctl
ADD loop.sh /loop
RUN chmod +x /loop

CMD ["/loop"]
