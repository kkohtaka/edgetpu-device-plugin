default: build

.PHONY: build
build:
	docker build -t quay.io/kkohtaka/edgetpu-demo:amd64 .

.PHONY: build-arm32
build-arm32:
	docker build -t quay.io/kkohtaka/edgetpu-demo:arm32 \
		--build-arg DEBIAN_BASE_SUFFIX=arm \
		.

.PHONY: run
run:
	kubectl apply -f edgetpu-demo.yaml

.PHONY: clean
clean:
	kubectl delete --ignore-not-found -f edgetpu-demo.yaml
