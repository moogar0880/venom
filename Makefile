# TEST_RESULTS defines the directory to which test results will be saved.
TEST_RESULTS=
LINT_RESULTS=

.PHONY: gogen
gogen:
	@go run generate_coercers.go
	@go run generate_getters.go

.PHONY: lint
lint:
ifeq ($(strip $(TEST_RESULTS)),)
	gofmt -d -l -s *.go
else
	mkdir -p $(LINT_RESULTS)
	gofmt -d -l -s *.go > $(LINT_RESULTS)/linter.out
endif

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint fmt
	@golangci-lint run

.PHONY: test
test:
ifeq ($(strip $(TEST_RESULTS)),)
	go test -v -coverprofile=coverage.out
else
	mkdir -p $(TEST_RESULTS)
	gotestsum --junitfile ${TEST_RESULTS}/gotestsum-report.xml -- $(go list ./... | circleci tests split --split-by=timings --timings-type=classname)
endif

.PHONY: test/benchmark
test/benchmark:
	go test -run=XXX -bench=. -benchmem

.PHONY: test/coverage
test/coverage: test
	go tool cover -html=coverage.out

.PHONY: vendor
vendor:
	go mod vendor
