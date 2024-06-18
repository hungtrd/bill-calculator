package main

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	// ...
	token := os.Getenv("openaikey")

	// init openai client
	openai.NewClient(token)

	// send request to openai
	// print response
}
