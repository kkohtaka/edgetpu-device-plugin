ARG DEBIAN_BASE_SUFFIX=amd64
FROM k8s.gcr.io/debian-base-${DEBIAN_BASE_SUFFIX}:0.4.1 as builder
ARG GO_TARBALL=go1.11.linux-amd64.tar.gz
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN apt-get update && apt-get install -y \
        curl \
    && rm -rf /var/lib/apt/lists/*
RUN curl https://dl.google.com/go/${GO_TARBALL} | tar zxv -C /usr/local
WORKDIR $GOPATH/src/github.com/revoman/edgetpu-device-plugin
COPY main.go .
COPY pkg pkg
COPY vendor vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/edgetpu-device-plugin
RUN curl http://storage.googleapis.com/cloud-iot-edge-pretrained-models/edgetpu_api.tar.gz | tar xzv -C /

FROM k8s.gcr.io/debian-base-${DEBIAN_BASE_SUFFIX}:0.4.1
ARG SO_SUFFIX=x86_64
ARG LIB_PATH=/lib/x86_64-linux-gnu
COPY --from=builder /python-tflite-source/libedgetpu/libedgetpu_${SO_SUFFIX}.so ${LIB_PATH}/libedgetpu.so
COPY --from=builder /python-tflite-source/99-edgetpu-accelerator.rules /etc/udev/rules.d/
COPY --from=builder /bin/edgetpu-device-plugin /bin/
RUN apt-get update && apt-get install -y \
        libusb-1.0 \
        udev \
    && rm -rf /var/lib/apt/lists/* \
    && udevadm trigger

RUN chgrp -R 0 /python-tflite-source \
    && chmod -R g=u /python-tflite-source
ENTRYPOINT ["/bin/edgetpu-device-plugin"]
USER 1001

