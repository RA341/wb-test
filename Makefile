TAG ?= wb

dk:
	@echo "Building Docker image with tag: $(TAG)"
	docker build . -t $(TAG)
