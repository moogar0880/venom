.PHONY: test
test:
	go test -coverprofile=coverage.out

.PHONY: vendor
vendor:
	dep ensure
