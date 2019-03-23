FROM k8s.gcr.io/debian-base-amd64:0.4.1 as builder
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get update && apt-get install -y \
        curl \
        build-essential \
        pkg-config \
        libusb-dev \
        libusb-1.0 \
    && rm -rf /var/lib/apt/lists/*
RUN curl https://dl.google.com/go/go1.11.linux-amd64.tar.gz | tar zxv -C /usr/local
WORKDIR $GOPATH/src/github.com/kkohtaka/edgetpu-device-plugin
COPY main.go .
COPY pkg pkg
COPY vendor vendor
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /bin/edgetpu-device-plugin
RUN curl http://storage.googleapis.com/cloud-iot-edge-pretrained-models/edgetpu_api.tar.gz | tar xzv -C /

FROM k8s.gcr.io/debian-base-amd64:0.4.1
COPY --from=builder /python-tflite-source/libedgetpu/libedgetpu_x86_64.so /lib/x86_64-linux-gnu/libedgetpu.so
COPY --from=builder /python-tflite-source/99-edgetpu-accelerator.rules /etc/udev/rules.d/
COPY --from=builder /bin/edgetpu-device-plugin /bin/
RUN apt-get update && apt-get install -y \
        libusb-1.0 \
        udev \
    && rm -rf /var/lib/apt/lists/* \
    && udevadm trigger
ENTRYPOINT ["/bin/edgetpu-device-plugin"]
