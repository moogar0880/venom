.PHONY: gogen
gogen:
	@go run generate_coercers.go
	@go run generate_getters.go

.PHONY: test
test:
	go test -coverprofile=coverage.out

.PHONY: test/coverage
test/coverage: test
	go tool cover -html=coverage.out

.PHONY: vendor
vendor:
	dep ensure
