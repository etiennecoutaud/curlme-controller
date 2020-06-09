GOOPTS=GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
DOCKER_IMG="etiennecoutaud/curlme-controller"
DOCKER_TAG="latest"

.PHONY: all build run-local test lint docker docker-push

all: lint test docker push

build:
	${GOOPTS} go build -o curlme-controller cmd/main.go

run-local: build
	./curlme-controller --kubeconfig=${KUBECONFIG}

run-docker:
	docker run ${DOCKER_IMG}:${DOCKER_TAG}

test:
	 go test -v -cover ./internal/.../

lint:
	golint -set_exit_status  ./...

docker:
	docker build -t ${DOCKER_IMG}:${DOCKER_TAG} .

docker-push:
	docker push ${DOCKER_IMG}:${DOCKER_TAG}