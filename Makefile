TAG ?= wb

dk:
	@echo "Building Docker image with tag: $(TAG)"
	docker build . -t $(TAG)
	docker run --rm -p 8022:8080 $(TAG)
