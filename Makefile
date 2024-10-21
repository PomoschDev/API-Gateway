.PHONY: all

ENTRY_POINT="./cmd/entrypoint"
EXE_NAME="./apigateway.exe"
ARGS ?= "--config=./config/local.yaml"
DOCKER_DIR="docker"

all: build remove_log run

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

remove_log:
	@echo Delete dir with log
	@del /s /q log -y
	@rmdir log