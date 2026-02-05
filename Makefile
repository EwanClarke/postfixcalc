.PHONY: test
test:
	go test ./internal/evaluator ./internal/parser ./internal/lexer

.PHONY: test-verbose
test-verbose:
	go test -v ./internal/evaluator ./internal/parser ./internal/lexer

.PHONY: test-coverage
test-coverage:
	go test -cover ./internal/evaluator ./internal/parser ./internal/lexer