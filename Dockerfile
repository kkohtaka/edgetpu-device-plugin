ARG DEBIAN_BASE_SUFFIX=amd64
FROM debian:bullseye-slim as builder
ARG GO_TARBALL=go1.17.linux-amd64.tar.gz
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get update && apt-get install -y \
        curl \
        git \
    && rm -rf /var/lib/apt/lists/*
RUN curl https://dl.google.com/go/${GO_TARBALL} | tar zxv -C /usr/local
WORKDIR $GOPATH/src/github.com/therevoman/edgetpu-device-plugin
COPY vendor vendor
COPY main.go .
COPY pkg pkg
COPY go.sum .
COPY go.mod .
#RUN go mod vendor \
#    && go get -u golang.org/x/xerrors \
#    && go get -u k8s.io/client-go@v0.18.19
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/edgetpu-device-plugin

FROM debian:bullseye-slim
ARG SO_SUFFIX=x86_64
ARG LIB_PATH=/lib/x86_64-linux-gnu
RUN apt update && apt install -y \
        curl \
        libusb-dev \
        udev \
        gnupg \
    && rm -rf /var/lib/apt/lists/* \
    && echo "deb https://packages.cloud.google.com/apt coral-edgetpu-stable main" | tee /etc/apt/sources.list.d/coral-edgetpu.list \
    && echo "deb https://packages.cloud.google.com/apt coral-cloud-stable main" | tee /etc/apt/sources.list.d/coral-cloud.list \
    && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
    && apt update && apt upgrade -y\
    && apt install -y \
        python3-pycoral \
        libedgetpu1-std \
        gasket-dkms \
        edgetpu-compiler \
    && rm -rf /var/lib/apt/lists/* \
    && (udevadm trigger || true)

ENTRYPOINT ["/bin/edgetpu-device-plugin"]
