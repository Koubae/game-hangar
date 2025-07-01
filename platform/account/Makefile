.PHONY: run build stop tests

# ============================
# 	Run
# ============================

run-reload:
	@air

run-local:
	go run cmd/main.go

# ============================
# 	Tests
# ============================

### Using locally installed Go
test-unit-local:
	go test -v ./... -coverprofile=coverage.out && go tool cover -o coverage.html -html=coverage.out
