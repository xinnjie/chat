ko-apply:
	kubectl kustomize config | KO_DOCKER_REPO=harbor.xinnjiedev.com/tennisapp ko apply -f -

proto-gen:
	buf generate

build:
	go build -o .bin ./server

test:
	go test ./...

mockgen:
	mockgen -source=server/auth/auth.go -destination=server/auth/mock_auth/mock_auth.go
	mockgen -source=server/store/store.go -destination=server/store/mock_store/mock_store.go

build-swift:
	swift build