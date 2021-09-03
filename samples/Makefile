REGISTRY=quay.io
PRODUCT=revoman/edgetpu-demo

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
	kubectl delete job edgetpu-demo --ignore-not-found
	kubectl delete pod -l job-name=edgetpu-demo --ignore-not-found
	kubectl apply -f edgetpu-demo.amd64.yaml

.PHONY: clean
clean:
	kubectl delete --ignore-not-found -f edgetpu-demo.amd64.yaml

.PHONY: push
push:
	docker push $(REGISTRY)/$(PRODUCT):amd64

.PHONY: push-arm32
push-arm32:
	docker push $(REGISTRY)/$(PRODUCT):arm32
