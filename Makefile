DOCKER_BUILDKIT=1
PLATFORM=local

all:
	docker build \
		--target bin-unix \
		--platform ${PLATFORM} \
		-t rdepot-cli:local \
		.

