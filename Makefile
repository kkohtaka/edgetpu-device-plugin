REGISTRY=quay.io
PRODUCT=kkohtaka/edgetpu-device-plugin

default: build

all: clean build samples

.PHONY: build
build:
	docker build -t $(REGISTRY)/$(PRODUCT):amd64 .

.PHONY: build-arm32
build-arm32:
	docker build -t $(REGISTRY)/$(PRODUCT):arm32 \
		--build-arg DEBIAN_BASE_SUFFIX=arm \
		--build-arg GO_TARBALL=go1.11.linux-armv6l.tar.gz \
		--build-arg SO_SUFFIX=arm32 \
		--build-arg LIB_PATH=/lib/arm-linux-gnueabihf \
		.

.PHONY: samples
samples:
	$(MAKE) -C samples

.PHONY: samples-arm32
samples-arm32:
	$(MAKE) -C samples build-arm32

.PHONY: clean
clean:
	$(MAKE) -C samples clean
	kubectl delete --ignore-not-found -f edgetpu-device-plugin.yaml

.PHONY: push
push:
	docker push $(REGISTRY)/$(PRODUCT):amd64
	$(MAKE) -C samples push

.PHONY: push-arm32
push-arm32:
	docker push $(REGISTRY)/$(PRODUCT):arm32
	$(MAKE) -C samples push-arm32
