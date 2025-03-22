ko-apply:
	kubectl kustomize config | KO_DOCKER_REPO=harbor.xinnjiedev.com/tennisapp ko apply -f -

proto-gen:
	buf generate

build:
	go build -o .bin ./server