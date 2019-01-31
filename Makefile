# Variables
PROJECT_NAMESPACE:=$(if $(PROJECT_NAMESPACE),$(PROJECT_NAMESPACE),$(shell cd .. && basename `pwd`))
PROJECT_NAME:=$(if $(PROJECT_NAME),$(PROJECT_NAME),$(shell basename `pwd`))
OUTPUT?=$(PROJECT_NAME)

test:
	@go test ./...

test-integration:
	@sh -c "source test/integration/integration.sh"

ensure:
	@go mod download

.PHONY: test test-integration
