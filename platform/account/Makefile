.PHONY: run build stop tests

# ============================
# 	Docker
# ============================

up-mongo-db:
	 @docker compose up db-mongodb-dashboard
down-mongo-db:
	 @docker compose down db-mongodb-dashboard mongo


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

test-specific:
ifndef TEST
	@echo "Please provide a test pattern using TEST=<pattern>"
	@echo "Example: make test-specific TEST=TestGetEnv/string_tests"
	@echo "make test-specific TEST=TestGetEnv"
	@echo "make test-specific TEST=TestGetEnv/string"
	@echo "make test-specific TEST=TestGetEnv/int"
	@echo "make test-specific TEST=TestGetEnv/int"
	@echo "make test-specific TEST=TestGetEnv/int_tests"
	@echo "\nAvailable test patterns:"
	@go test ./... -v -list=. | grep "^Test"
else
	go test ./... -v -run $(TEST)
endif
