FROM k8s.gcr.io/debian-base-amd64:0.4.1
RUN apt-get update \
    && apt-get install --yes \
        sudo \
        bash \
        curl \
        ca-certificates \
        build-essential \
    && mkdir -p /etc/udev/rules.d/ \
    && curl http://storage.googleapis.com/cloud-iot-edge-pretrained-models/edgetpu_api.tar.gz | tar xzv -C / \
    && cd /python-tflite-source \
    && yes | bash ./install.sh \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /python-tflite-source/edgetpu/
ENTRYPOINT ["python3.5"]
CMD ["--version"]
