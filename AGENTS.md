# Guidelines for LLM Agents

- Use `go.uber.org/mock` for generating mocks
- Use `github.com/stretchr/testify` for assertions in tests
- MUST run `go test -race ./...` and `golangci-lint run ./...` before finalizing any code changes