ARG DEBIAN_BASE_SUFFIX=amd64
FROM debian:bullseye-slim
ARG DEBIAN_BASE_SUFFIX
USER root

ADD platform_recognizer.${DEBIAN_BASE_SUFFIX}.sh /tmp/
RUN apt update \
    && apt install --yes \
        sudo \
        bash \
        curl \
        gnupg \
    && echo "deb https://packages.cloud.google.com/apt coral-edgetpu-stable main" | tee /etc/apt/sources.list.d/coral-edgetpu.list \
    && echo "deb https://packages.cloud.google.com/apt coral-cloud-stable main" | tee /etc/apt/sources.list.d/coral-cloud.list \
    && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add - \
    && mkdir -p /etc/udev/rules.d/ \
    && mkdir /python-tflite-source \
    && mv /tmp/platform_recognizer.${DEBIAN_BASE_SUFFIX}.sh /python-tflite-source/platform_recognizer.sh \
    && apt update && apt upgrade --yes \
    && apt install --yes \
        ca-certificates \
        build-essential \
        git \
        edgetpu-compiler \
        libedgetpu1-std \
        libedgetpu-dev \
        gasket-dkms \
        python3-pycoral \
        python3-tflite-runtime \
    && cd /python-tflite-source \
    && git clone https://github.com/google-coral/pycoral.git \
    && cd pycoral \
    && yes | bash examples/install_requirements.sh classify_image.py \
    && rm -rf /var/lib/apt/lists/*

RUN chgrp -R 0 /python-tflite-source \
    && chmod -R g=u /python-tflite-source

USER 1001
WORKDIR /python-tflite-source/pycoral/
ENTRYPOINT ["python3"]
CMD ["--version"]
