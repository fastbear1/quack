Go ?= go
SHELL=/bin/bash
TODAY=$(shell date +'%Y-%m-%d 00:00')

.PHONY: tests version

test:
	@cp ./misc/quack_config_test.yaml ./quack_config.yaml
	@go test -v ./cmd ./internal ./drivers ./runner
	@cp ./misc/quack_config_original.yaml ./quack_config.yaml

version:
	@MAJOR=0; \
	MINOR=$$(git log --date=short --pretty=format:%ad | sort | uniq -c | wc -l); \
	PATCH=$$(git log --date=short --pretty=format:%ad --after='$(TODAY)' |sort|uniq -c|awk {'print $$1 + 1'}); \
	if [ "$$PATCH" == "" ]; then PATCH=1; fi; \
	echo "$$MAJOR"."$$MINOR"."$$PATCH"
