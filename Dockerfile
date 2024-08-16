FROM ubuntu:20.04

RUN apt update -y
RUN apt install wget tar gzip git -y

RUN DEBIAN_FRONTEND=noninteractive apt install build-essential libcap-dev pkg-config libsystemd-dev -y
RUN wget -P /tmp https://mirror.ghproxy.com/https://github.com/ioi/isolate/archive/master.tar.gz && tar -xzvf /tmp/master.tar.gz -C / > /dev/null
RUN make -C /isolate-master install && rm /tmp/master.tar.gz

# python 
RUN apt install python3.8 -y
RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 1
# golang
WORKDIR /usr/local
ARG GO_VERSION=1.22.4
ARG GO_ARCH=linux-amd64
RUN wget -q -o /dev/null https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz
RUN tar -xzf go${GO_VERSION}.${GO_ARCH}.tar.gz
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go
ENV GOPROXY=https://goproxy.cn
ENV PATH=${GOROOT}/bin:${GOPATH}/bin:${PATH}
# java
WORKDIR /tmp
RUN wget https://download.oracle.com/java/17/latest/jdk-17_linux-x64_bin.tar.gz && tar -xzvf ./jdk-17_linux-x64_bin.tar.gz  > /dev/null
RUN mkdir -p /usr/local/jdk17
RUN mv ./jdk-17*/* /usr/local/jdk17 && rm ./jdk-17_linux-x64_bin.tar.gz 
ENV JAVA_HOME=/usr/local/jdk17
ENV PATH=$JAVA_HOME/bin:$PATH

WORKDIR /sandbox

COPY ./ .
WORKDIR /sandbox/cmd
RUN go mod download
RUN go build -o ./server .
ARG SERVICE_PORT
ENV SERVICE_PORT=${SERVICE_PORT}
EXPOSE ${SERVICE_PORT}

CMD ["./server"]

