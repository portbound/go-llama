# Makefile

BINARY_NAME=llama-repl
CMD_PATH=./cmd/llama-repl

build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

run: build
	# Start the Ollama server in the background
	ollama serve &

	# Open a new Windows Terminal window running the REPL
	wt.exe wsl -e ./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

