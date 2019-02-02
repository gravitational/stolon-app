FROM quay.io/gravitational/debian-venti:go1.11.5-stretch as builder
ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm \
    REPO_PATH=/gopath/src/github.com/gravitational/stolon
WORKDIR $REPO_PATH
COPY stolon .

RUN mkdir -p bin && \
    CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o bin/stolon-sentinel ${REPO_PATH}/cmd/sentinel && \
    CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o bin/stolon-proxy ${REPO_PATH}/cmd/proxy && \
    go build -a -installsuffix cgo -ldflags '-s' -o bin/stolon-keeper ${REPO_PATH}/cmd/keeper

FROM quay.io/gravitational/debian-grande:stretch

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        postgresql-9.6 \
        postgresql-client-9.6 \
        postgresql-contrib-9.6

COPY --from=builder /gopath/src/github.com/gravitational/stolon/bin/ /usr/local/bin/
ADD rootfs/ /

RUN chmod a+rx /usr/local/bin/run.sh
RUN useradd -ms /bin/bash stolon

EXPOSE 5431 5432 6431

RUN test -f /cleanup.sh && sh /cleanup.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/run.sh"]
