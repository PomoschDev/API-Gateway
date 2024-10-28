.PHONY: all

ENTRY_POINT="./cmd/entrypoint"
EXE_NAME="./apigateway.exe"
ARGS ?= "--config=./config/local.yaml"
DOCKER_DIR="docker"

all: build run

build:
	@echo Build and compile
	@go build -o ${EXE_NAME} ${ENTRY_POINT}
	@echo Good build

run:
	@echo Run ${EXE_NAME} ${ARGS}
	@${EXE_NAME} ${ARGS}

compose_docker:
	@echo Docker
	@docker compose up -d

compose_clean:
	@echo Docker without repository
	@cd ${DOCKER_DIR} && docker compose up -d

doc:
	@echo Make swagger docs
	@swag init