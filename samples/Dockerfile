ARG DEBIAN_BASE_SUFFIX=amd64
FROM k8s.gcr.io/debian-base-${DEBIAN_BASE_SUFFIX}:0.4.1
ARG DEBIAN_BASE_SUFFIX
ADD platform_recognizer.${DEBIAN_BASE_SUFFIX}.sh /tmp/
RUN apt-get update \
    && apt-get install --yes \
        sudo \
        bash \
        curl \
        ca-certificates \
        build-essential \
    && mkdir -p /etc/udev/rules.d/ \
    && curl http://storage.googleapis.com/cloud-iot-edge-pretrained-models/edgetpu_api.tar.gz | tar xzv -C / \
    && mv /tmp/platform_recognizer.${ARCH}.sh /python-tflite-source/platform_recognizer.sh \
    && cd /python-tflite-source \
    && yes | bash ./install.sh \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /python-tflite-source/edgetpu/
ENTRYPOINT ["python3.5"]
CMD ["--version"]
