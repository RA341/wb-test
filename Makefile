TAG ?= wb

dk:
	@echo "Building Docker image with tag: $(TAG)"
	docker build . -t $(TAG)
	docker run --rm -p 9933:9933 $(TAG)
