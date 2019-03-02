.PHONY: gogen
gogen:
	@go run generate_coercers.go
	@go run generate_getters.go

.PHONY: lint
lint:
	gofmt -d -l -s *.go

.PHONY: test
test:
	go test -coverprofile=coverage.out

.PHONY: test/benchmark
test/benchmark:
	go test -run=XXX -bench=. -benchmem

.PHONY: test/coverage
test/coverage: test
	go tool cover -html=coverage.out

.PHONY: vendor
vendor:
	go mod vendor
