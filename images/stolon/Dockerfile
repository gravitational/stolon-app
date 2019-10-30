ARG PGVERSION=9.6
# base build image
FROM golang:1.13-buster AS build_base
ARG STOLON_REVISION=656bc2fe5d801193b412bcfb57dbfa84cfd9f386

RUN set -eux; \
    apt-get update && \
    apt-get install -qq -y git && \
    git clone https://github.com/gravitational/stolon.git /stolon && \
    cd /stolon && git checkout $STOLON_REVISION

WORKDIR /stolon

RUN go mod download

#######
####### Build the stolon binaries
#######
FROM build_base AS builder

RUN make

#######
####### Build the final image
#######
FROM postgres:$PGVERSION

RUN useradd -ms /bin/bash stolon

EXPOSE 5432

COPY --from=builder /stolon/bin/ /usr/local/bin

RUN chmod +x /usr/local/bin/stolon-keeper /usr/local/bin/stolon-sentinel /usr/local/bin/stolon-proxy /usr/local/bin/stolonctl
