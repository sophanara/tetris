build:
	@go build -o bin/tetris cmd/main.go

run: build
	@./bin/tetris


