RUBY_VERSION=3.1
DOCKER_IMAGE=otel-ruby
DOCKER_TAG=latest

.PHONY: all
all:
	@$(MAKE) build-image
	@$(MAKE) mount-container
	@$(MAKE) copy-files
	@$(MAKE) unmount-container

build-image:
	@echo "\n🚀 Building image"
	@docker build \
		--build-arg RUBY_VERSION=$(RUBY_VERSION) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-f Dockerfile .

mount-container:
	@echo "🧪 Mounting container"
	@docker create --name $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(DOCKER_TAG)

unmount-container:
	@echo "🧪 Unmounting container"
	@docker rm $(DOCKER_IMAGE)

copy-files:
	@echo "🧪 Copying files"
	@docker cp $(DOCKER_IMAGE):/Gemfile.lock ./Gemfile.lock
