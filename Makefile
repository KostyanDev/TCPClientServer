include .env
export

BIN_DIR = bin/
CLIENT_MAIN = cmd/client/main.go
CLIENT_NAME = client
SERVER_NAME = tcp_server
SERVER_IMAGE_NAME = tcp_server

default: server-build server-run client-build client-run

client-build: | $(BIN_DIR)
	@printf "\033[32;1mBuild client\033[0m\n"
	@go build -o $(BIN_DIR)$(CLIENT_NAME) $(CLIENT_MAIN)

client-run:
	@printf "\033[32;1mRun client\033[0m\n"
	@$(BIN_DIR)$(CLIENT_NAME)

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

server-run:
	@printf "\033[32;1mRun server\033[0m\n"
	@docker run -d --rm -e SERVER_PORT=$(PORT) -p $(PORT):$(PORT) --name $(SERVER_NAME) $(SERVER_IMAGE_NAME)

server-build:
	@printf "\033[32;1mBuild server\033[0m\n"
	@docker build -t $(SERVER_IMAGE_NAME) .

clean:
	@rm -rf $(BIN_DIR)
	@docker stop $(SERVER_NAME)
	@docker rmi $(SERVER_IMAGE_NAME)
	@docker image prune

.PHONY: server-run server-build default
