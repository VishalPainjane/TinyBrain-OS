.PHONY: build-image dev-shell test-cgo

# Build the development container
build-image:
	docker build -t tinybrain-dev -f Dockerfile.dev .

# Drop into an interactive shell inside the CUDA-enabled container
dev-shell: build-image
	docker run -it -v $(CURDIR):/app tinybrain-dev /bin/bash

# Shortcut to run the CGO tests safely inside the Linux container
test-cgo: build-image
	docker run -v $(CURDIR):/app tinybrain-dev /bin/bash -c "cd internal/scheduler/v2 && nvcc -c backend.cu -o backend.o -O3 -lcublas && ar rcs libtinybrain.a backend.o && cd ../../../ && go test -v ./internal/scheduler/v2/..."
