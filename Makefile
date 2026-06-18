.PHONY: tests

tests:
	@cp ./misc/quack_config_test.yaml ./quack_config.yaml
	@go test -v ./cmd ./internal ./drivers
	@cp ./misc/quack_config_original.yaml ./quack_config.yaml
