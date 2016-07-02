FROM alpine:3.4

ADD bin/stolonctl /usr/local/bin/stolonctl
ADD loop.sh /loop
RUN chmod +x /loop

CMD ["/loop"]
