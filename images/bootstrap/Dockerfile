FROM ubuntu:14.04

RUN apt-get update && apt-get install -y curl python2.7

RUN curl -O https://storage.googleapis.com/kubernetes-release/release/v1.2.0/bin/linux/amd64/kubectl && \
    cp kubectl /usr/local/bin && \
    chmod +x /usr/local/bin/kubectl

ADD resources /resources
ADD bootstrap.py /bootstrap.py

CMD ["/usr/bin/python", "/bootstrap.py"]
