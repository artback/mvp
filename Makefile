all : test  vet

test:
	go test ./...

test-integration:
	go test ./... -tags=integration

vet :
	go vet ./...

