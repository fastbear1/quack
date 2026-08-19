Go ?= go
SHELL=/bin/bash
TODAY=$(shell date +'%Y-%m-%d 00:00')

.PHONY: all

test:
	@cp ./misc/quack_config_test.yaml ./quack_config.yaml
	@go test -v ./cmd ./internal ./drivers/pg_driver ./runner
	@cp ./misc/quack_config_original.yaml ./quack_config.yaml

int-test:
	@docker-compose -f misc/docker-compose.yaml up -d
	@PGPASSWORD=pass psql -U quack -h localhost -d quack -p 5432 -c "DROP TABLE IF EXISTS test_table_c;DROP TABLE IF EXISTS test_table_b;DROP TABLE IF EXISTS test_table_a;"
	@rm -rf tests/migrations
	@mkdir tests/migrations
	@go test -v ./tests

version:
	@MAJOR=0; \
	MINOR=$$(git log --date=short --pretty=format:%ad | sort | uniq -c | wc -l); \
	PATCH=$$(git log --date=short --pretty=format:%ad --after='$(TODAY)' |sort|uniq -c|awk {'print $$1 + 1'}); \
	if [ "$$PATCH" == "" ]; then PATCH=1; fi; \
	echo "$$MAJOR"."$$MINOR"."$$PATCH"

build:
	@go build -ldflags="-s -w" -o quack cmd/main.go

help:
	@echo "Command usage:"
	@echo "     help - show current help information"
	@echo "     test - run all tests"
	@echo "     version - show application version for uncommited changes"

