FROM dunedaq/sl7-minimal
LABEL maintainer="Glenn Dirkx <glenn.dirkx@cern.ch>"

RUN yum install -y make && \
    yum clean all

# golang
ENV GO_VERSION=1.15.3
ENV PATH="$PATH:/usr/local/go/bin"
RUN curl -o go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && rm -f go.tar.gz