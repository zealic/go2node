# Variables
PROJECT_NAMESPACE:=$(if $(PROJECT_NAMESPACE),$(PROJECT_NAMESPACE),$(shell cd .. && basename `pwd`))
PROJECT_NAME:=$(if $(PROJECT_NAME),$(PROJECT_NAME),$(shell basename `pwd`))
OUTPUT?=$(PROJECT_NAME)

test:
	@go test ./...
	node channel_parent_test.js channel_child

ensure:
	@go mod download
