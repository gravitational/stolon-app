FROM stolon:0.2.0

ADD /backup.sh /backup
ADD /restore.sh /restore
RUN chmod +x /backup /restore

ENTRYPOINT ["/backup"]
