.PHONY: run build stop tests


quickstart: init postgres-up migrate-identity-up
	@echo "Quickstart completed successfully"


# ============================
# 	Run
# ============================
# //////////////////////
# 	local
# //////////////////////
run-identity-local:
	@air -c .air.identity.toml

run-identity-local-no-hot-reload:
	@go run ./cmd/identity/main.go


# //////////////////////
# 	docker
# //////////////////////

# ······················
#		DB -- PostGreSQL
# ······················
postgres-up:
	@docker compose up db-postgres-dashboard
postgres-down:
	@docker compose down db-postgres-dashboard
postgres-down-clean-up:
	@docker compose down -v db-postgres db-postgres-dashboard

postgres-shell:
	docker compose exec \
		-e PGPASSWORD='admin' \
		db-postgres \
		psql -U admin -d game_hangar

# ============================
# 	Init
# ============================
init: install-deps update-env-file

install-deps:
	go mod tidy

update-env-file:
	@echo 'Updating .env from .env.example 🖋️...'
	@cp .env.example .env
	@echo '.env Updated ✨'

# ============================
# 	Tests
# ============================
COVERAGE_THRESHOLD ?= 80
COVERAGE_PKGS := $(shell go list ./... | grep -v '/internal/mocks' | grep -v '/pkg/generated' | grep -v '/cmd/demo' | grep -v '/cmd' | grep -v '/internal/run')

test-all:
	go test -v $(COVERAGE_PKGS) -cover


test-unit:
	go test -v -short $(COVERAGE_PKGS) -cover

# TODO: Check whether there is a better way to do this. This was AI generated and seems a mess
# Intention here is:
#	1) Ignore certain folders (COVERAGE_PKGS should list ONLY the actual go module that are testable)
#	2) Spit an coverage.out and translate to html
#	3) Grabs the total and then checks the coverage threshold is greater or equal the threshold)
test-all-coverage-html:
	@go test -v $(COVERAGE_PKGS) -coverprofile=coverage.out && \
	go tool cover -o coverage.html -html=coverage.out && \
	coverage=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	awk "BEGIN { if ($$coverage < $(COVERAGE_THRESHOLD)) { printf \"Coverage check failed: %.1f%% < %s%%\n\", $$coverage, \"$(COVERAGE_THRESHOLD)\"; exit 1 } else { printf \"Coverage check passed: %.1f%% >= %s%%\n\", $$coverage, \"$(COVERAGE_THRESHOLD)\" } }"


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
	@go test ./... -v -run $(TEST)
endif

# ============================
# 	Scripts
# ============================


# //////////////////////
# 	Migrations
# //////////////////////
migrate-identity-up:
	@go run ./migrations/identity/migrate_identity.go -action up 0
migrate-identity-down:
	@go run ./migrations/identity/migrate_identity.go -action down 0
migrate-identity-status:
	@go run ./migrations/identity/migrate_identity.go -action status
