GOBIN=$(shell pwd)/bin

gethulent:
	GOBIN=$(GOBIN) go install cmd/gethulent.go

test:
	go test -cover ./...

nice:
	go fmt ./...

.PHONY: gethulent test nice