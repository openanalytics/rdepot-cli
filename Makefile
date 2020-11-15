DOCKER_BUILDKIT=1
PLATFORM=local

all:
	docker build \
		--target bin-unix \
		--platform ${PLATFORM} \
		-t rdepot-cli:local \
		.

license:
	# go get -u github.com/google/addlicense
	addlicense \
		-c "Open Analytics" \
		-y 2020 \
		-l apache \
		*.go */*.go

