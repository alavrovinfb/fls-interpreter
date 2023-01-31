PROJECT_ROOT            ?= github.com/alavrovinfb/fls-interpreter
#$(PWD)
#github.com/fls-interpreter
BUILD_PATH              ?= bin
DOCKERFILE_PATH         := ./docker
DOCKER_FILE			    := $(DOCKERFILE_PATH)/Dockerfile

# configuration for the protobuf gentool
SRCROOT_ON_HOST         := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_IN_CONTAINER    := /go/src/$(PROJECT_ROOT)
DOCKER_RUNNER           ?= docker run --rm -u `id -u`:`id -g` -e GOCACHE=/go -e CGO_ENABLED=0
DOCKER_RUNNER           += -v $(SRCROOT_ON_HOST):$(SRCROOT_IN_CONTAINER)
DOCKER_GENERATOR        := infoblox/atlas-gentool:v27
GENERATOR               := $(DOCKER_RUNNER) $(DOCKER_GENERATOR)

PROTOBUF_ARGS =	 -I=$(PROJECT_ROOT)/vendor
PROTOBUF_ARGS += --go_out=.
PROTOBUF_ARGS += --go-grpc_out=.
PROTOBUF_ARGS += --go-grpc_opt require_unimplemented_servers=false
PROTOBUF_ARGS += --validate_out="lang=go:."
PROTOBUF_ARGS += --grpc-gateway_out=.
PROTOBUF_ARGS += --grpc-gateway_opt logtostderr=true,allow_delete_body=true
PROTOBUF_ARGS += --openapiv2_out=.
PROTOBUF_ARGS += --openapiv2_opt allow_delete_body=true,atlas_patch=true,json_names_for_fields=false

# configuration for image names
USERNAME                := $(USER)
SERVICE_NAME            := fls-interpreter
GIT_COMMIT              := $(shell git describe --long --tags --dirty=-unreleased --always || echo pre-commit)
IMAGE_VERSION           := $(GIT_COMMIT)-j$(BUILD_NUMBER)
BUILD_NUMBER        	?= 0
IMAGE_REGISTRY 			?= $(USERNAME)
IMAGE_NAME              ?= $(IMAGE_REGISTRY)/$(SERVICE_NAME)

.PHONY test:
test:
	@go test ./... -cover
.PHONY build-local:
build-local:
	@go build -o bin/fls-interpreter ./cmd

.PHONY run-local:
run-local:
	@go run ./cmd --script.files="example/sample-script.txt"

.PHONY build-docker: test
build-docker:
	@docker build --build-arg REPO="${PROJECT_ROOT}" \
		-f $(DOCKER_FILE) \
		-t $(IMAGE_NAME):latest \
		-t $(IMAGE_NAME):$(IMAGE_VERSION) .
	@docker image prune -f --filter label=stage=build-intermediate

.PHONY vendor: protobuf
vendor:
	@go mod tidy
	@go mod vendor

.PHONY protobuf:
protobuf:
	$(GENERATOR) \
	$(PROTOBUF_ARGS) \
	$(PROJECT_ROOT)/pkg/pb/service.proto
