# go-llama

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Ollama](https://img.shields.io/badge/Ollama-Required-brightgreen)](https://ollama.com)

**`go-llama`** is a minimal terminal-based REPL interface for interacting with an [Ollama](https://ollama.com)-compatible LLM API. It features chat history persistence, model selection, and a form-based TUI powered by [Charmbracelet](https://github.com/charmbracelet) components.

> 🦙 Designed for use with Ollama running locally (`http://localhost:11434`)

---

## Features

- 🧵 Multi-turn chat sessions with history
- 💾 Automatic chat log saving in `chats/`
- 🔄 Resume previous chats or start new ones
- 🧩 Interactive model selector from installed Ollama models
- 🎛️ Clean terminal UI built with [`huh`](https://github.com/charmbracelet/huh) and `bubbletea`

---

## Installation

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Ollama](https://ollama.com/) running locally with at least one model installed

### Clone and Run

```bash
git clone https://github.com/portbound/go-llama
cd go-llama
go run .
```

### Build a Binary

```bash
go build -o go-llama .
./go-llama
```
License

MIT © portbound
