default: build

all: clean build samples apply run-samples

.PHONY: build
build:
	docker build -t quay.io/kkohtaka/edgetpu-device-plugin:amd64 .

.PHONY: build-arm32
build-arm32:
	docker build -t quay.io/kkohtaka/edgetpu-device-plugin:arm32 \
		--build-arg DEBIAN_BASE_SUFFIX=arm \
		--build-arg GO_TARBALL=go1.11.linux-armv6l.tar.gz \
		--build-arg SO_SUFFIX=arm32 \
		--build-arg LIB_PATH=/lib/arm-linux-gnueabihf \
		.

.PHONY: apply
apply:
	kubectl apply -f edgetpu-device-plugin.yaml

.PHONY: samples
samples:
	$(MAKE) -C samples

.PHONY: samples-arm32
samples-arm32:
	$(MAKE) -C samples build-arm32

.PHONY: run-samples
run-samples:
	$(MAKE) -C samples run

.PHONY: clean-samples
clean-samples:
	$(MAKE) -C samples clean

.PHONY: clean
clean:
	$(MAKE) -C samples clean
	kubectl delete --ignore-not-found -f edgetpu-device-plugin.yaml

.PHONY: push
push:
	docker push quay.io/kkohtaka/edgetpu-device-plugin:amd64
	docker push quay.io/kkohtaka/edgetpu-demo:amd64

.PHONY: push-arm32
push-arm32:
	docker push quay.io/kkohtaka/edgetpu-device-plugin:arm32
	docker push quay.io/kkohtaka/edgetpu-demo:arm32
