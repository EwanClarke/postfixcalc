.PHONY: test
test:
	go test ./internal/engine

.PHONY: test-verbose
test-verbose:
	go test -v ./internal/engine

.PHONY: test-coverage
test-coverage:
	go test -cover ./internal/engine