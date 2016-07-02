FROM stolonctl:0.2.0

RUN apk update && \
    apk add curl

RUN curl -O https://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/linux/amd64/kubectl && \
    mv kubectl /usr/local/bin && \
    chmod +x /usr/local/bin/kubectl

ADD bin/stolonhatest /stolonhatest

CMD ["/stolonhatest"]
