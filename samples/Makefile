default: build

.PHONY: build
build:
	docker build -t edgetpu-demo:latest .

.PHONY: run
run:
	kubectl apply -f edgetpu-demo.yaml

.PHONY: clean
clean:
	kubectl delete --ignore-not-found -f edgetpu-demo.yaml
