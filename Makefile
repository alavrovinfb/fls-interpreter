PROJECT_ROOT            ?= $(PWD)
#github.com/fls-interpreter
BUILD_PATH              ?= bin
DOCKERFILE_PATH         := $(PROJECT_ROOT)/docker
DOCKER_FILE			    := $(DOCKERFILE_PATH)/Dockerfile

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
