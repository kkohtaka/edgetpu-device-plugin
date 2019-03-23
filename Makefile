default: build

all: clean build samples apply run-samples

.PHONY: build
build:
	docker build -t edgetpu-device-plugin:latest .

.PHONY: apply
apply:
	kubectl apply -f edgetpu-device-plugin.yaml

.PHONY: samples
samples:
	$(MAKE) -C samples

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
