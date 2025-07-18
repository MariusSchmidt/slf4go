GO = CGO_ENABLED=0 GO111MODULE=on go
DIR_REPORTS = ./reports

# font and color definitions
BOLD := $(shell tput bold)
RED := $(shell tput setaf 1)
GREEN := $(shell tput setaf 2)
YELLOW := $(shell tput setaf 3)
BLUE := $(shell tput setaf 4)
MAGENTA := $(shell tput setaf 5)
CYAN := $(shell tput setaf 6)
WHITE := $(shell tput setaf 7)
RESET := $(shell tput sgr0)


default: helper

helper: # Prints this help message
	@echo "Usage: \n"
	@grep '^[A-Za-z].*:' $(lastword $(MAKEFILE_LIST)) | sed 's/:.*#/#/' | sort | column -t -s'#'

clean: # Remove reports directory and cleans up go.mod file
	$(GO) mod tidy
	rm -rf $(DIR_REPORTS)

generate_mocks: # Generates all mock implementations
	@echo "\n$(BOLD)$(BLUE)📥 Install mock generators...$(RESET)"
	go get github.com/golang/mock/mockgen@latest
	@echo "$(GREEN)✓ Install complete$(RESET)"
	@echo "\n$(BOLD)$(BLUE)⚡ Generating mocks...$(RESET)"
	$(GO) generate ./...
	@echo "$(GREEN)✓ Generation complete$(RESET)"

run_static_checks: # Runs static code analysis
	@echo "\n$(BOLD)$(BLUE)📥 Install analysis tools...$(RESET)"
	go install golang.org/x/lint/golint@latest # Code style checker
	go install github.com/kisielk/errcheck@latest # Error handling checker
	go install honnef.co/go/tools/cmd/staticcheck@latest # Static code analyzer
	@echo "$(GREEN)✓ Install complete$(RESET)"
	@echo "\n$(BOLD)$(BLUE)🔍 Perform static checks...$(RESET)"
	@pwd && staticcheck -go module -f stylish ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Static checks complete$(RESET)"

run_tests: # Runs all available unit tests
	@echo "\n$(BOLD)$(BLUE)🔬 Running tests...$(RESET)"
	@mkdir -p $(DIR_REPORTS)
	$(GO) test -v -cover -coverpkg=$(shell go list ./... | grep -v test_mocks) -coverprofile=$(DIR_REPORTS)/coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=$(DIR_REPORTS)/coverage.out -o $(DIR_REPORTS)/coverage.html
	@echo "$(GREEN)✓ Tests completed$(RESET)"


show_version: # Displays the version of this module
	@cat version.txt && printf "\n"

all: clean generate_mocks run_static_checks run_tests # Runs all targets

.PHONY: helper clean generate_mocks run_static_checks run_tests show_version all

