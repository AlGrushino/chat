GOCMD=go
GOBUILD=$(GOCMD) build
NAME=chatb

all: clean build run

build:
	@$(GOBUILD) -o $(NAME) ./cmd 

clean:
	@rm -rf $(NAME) ./logs/app.log

run:
	@echo "Starting application (migrations skipped by default)..."
	@touch ./logs/app.log
	@./$(NAME)

migrate:
	@echo "Running migrations..."
	./$(NAME) --migrate

test:
	go test ./internal/service/chat/... -v