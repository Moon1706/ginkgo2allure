export GO111MODULE=on

.PHONY: bin
bin: fmt vet
	go build -o bin/ginkgo2allure github.com/Moon1706/ginkgo2allure

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test --short --count=1 --covermode=count --coverprofile=coverage.out --coverpkg=./... ./...

.PHONY: coverage
coverage: test
	go tool cover --func=coverage.out
