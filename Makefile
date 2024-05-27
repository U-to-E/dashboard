# Variables
BINARY_NAME=dashboard
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_MOD=$(GO_CMD) mod

# Colors
GREEN=\033[0;32m
RED=\033[0;31m
RESET=\033[0m

# Targets
.PHONY: all build run clean tidy docker-build docker-run

all: build

build:
	@echo "$(GREEN) Building binary...$(RESET)"
	$(GO_BUILD) -o $(BINARY_NAME)

run:
	@echo "$(GREEN)Running the application...$(RESET)"
	$(GO_RUN) main.go

clean:
	@echo "$(RED)Cleaning up...$(RESET)"
	rm -f $(BINARY_NAME)

tidy:
	@echo "$(GREEN)Tidying Go modules...$(RESET)"
	$(GO_MOD) tidy

docker-build:
	@echo "$(GREEN)Building Docker image...$(RESET)"
	docker build -t dashboard-app .

docker-run:
	@echo "$(GREEN)Running Docker container...$(RESET)"
	docker run --env-file .env -p 3000:3000 dashboard-app