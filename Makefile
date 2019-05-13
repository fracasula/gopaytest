IMAGE_NAME := gopaytest/payments-api
IMAGE_TAG := dev-latest
HTTP_PORT := 8080
BASE_URL := https://api.test.gopaytest.tech

default: build run

.PHONY: build
build:
	docker build \
		--build-arg HTTP_PORT=$(HTTP_PORT) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: run
run:
	docker run --rm -it \
		-e BASE_URL="$(BASE_URL)" \
		-p $(HTTP_PORT):$(HTTP_PORT) \
		$(IMAGE_NAME):$(IMAGE_TAG)
