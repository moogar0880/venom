.PHONY: gogen
gogen:
	@go run generate_coercers.go
	@go run generate_getters.go

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
	@if [[ $$(go version) = *"go1.11"* ]]; then\
        go mod vendor;\
	else\
		dep init;\
    fi
	
