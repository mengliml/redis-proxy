.PHONY: help

# HELP 
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# Build the container
build: ## Build the development container
	docker-compose build

run: ## Run container in development mode
	docker-compose up -d

stop: ## Stop running docker container
	docker-compose down -v

test: ## Run all tests 
	docker-compose run proxy go test -v ./...
	docker-compose down
