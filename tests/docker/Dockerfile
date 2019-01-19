FROM debian:jessie
ARG ANSIBLE_VERSION=2.7.4

ENV DEBIAN_FRONTEND=noninteractive
RUN echo "===> Installing python and supporting tools..." && \
    apt-get update -y && apt-get install --fix-missing && \
    apt-get install -y \
    python python-yaml \
    curl gcc python-pip python-dev libffi-dev libssl-dev && \
    apt-get -y --purge remove python-cffi && \
    pip install --upgrade setuptools && \
    pip install --upgrade pycrypto && \
    pip install --upgrade cffi && \
    pip install --upgrade requests google-auth && \
    \
    \
    echo "===> Fix strange bug under Jessie: cannot import name IncompleteRead" && \
    easy_install -U pip && \
    \
    \
    echo "===> Installing Ansible..." && \
    pip install ansible==${ANSIBLE_VERSION} && \
    echo "===> Removing unused APT resources..." && \
    apt-get -f -y --auto-remove remove \
    gcc python-pip python-dev libffi-dev libssl-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/*

# default command: display Ansible version
CMD [ "ansible-playbook", "--version" ]
