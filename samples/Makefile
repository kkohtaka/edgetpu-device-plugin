REGISTRY=quay.io
PRODUCT=kkohtaka/edgetpu-demo

default: build

.PHONY: build
build:
	docker build -t $(REGISTRY)/$(PRODUCT):amd64 .

.PHONY: build-arm32
build-arm32:
	docker build -t $(REGISTRY)/$(PRODUCT):arm32 \
		--build-arg DEBIAN_BASE_SUFFIX=arm \
		.

.PHONY: run
run:
	kubectl apply -f edgetpu-demo.yaml

.PHONY: clean
clean:
	kubectl delete --ignore-not-found -f edgetpu-demo.yaml

.PHONY: push
push:
	docker push $(REGISTRY)/$(PRODUCT):amd64

.PHONY: push-arm32
push-arm32:
	docker push $(REGISTRY)/$(PRODUCT):arm32
