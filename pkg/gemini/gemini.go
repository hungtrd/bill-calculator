package gemini

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	Model1_5Flash = "gemini-1.5-flash"
	Model1_5Pro   = "gemini-1.5-pro"
)

type GeminiClient struct {
	Client    *genai.Client
	ModelName string
}

func NewClient(model string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_TOKEN")))
	if err != nil {
		return nil, err
	}

	return &GeminiClient{Client: client, ModelName: model}, nil
}

func (g *GeminiClient) GenerateContentWithFile(ctx context.Context, filePath string, prompt string) (string, error) {
	startTime := time.Now()
	model := g.Client.GenerativeModel(g.ModelName)
	imgData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	promptData := []genai.Part{
		genai.ImageData("jpeg", imgData),
		genai.Text(prompt),
	}

	resp, err := model.GenerateContent(ctx, promptData...)
	if err != nil {
		return "", err
	}

	fmt.Println("total token usage: ", resp.UsageMetadata.TotalTokenCount)
	fmt.Println("total time taken: ", time.Since(startTime))

	if len(resp.Candidates) > 0 {
		can := resp.Candidates[0]
		if can.Content != nil && len(can.Content.Parts) > 0 {
			s := fmt.Sprintf("%v", can.Content.Parts[0])
			return s, nil
		}
	}

	return "", nil
}

func (g *GeminiClient) Close() {
	g.Client.Close()
}
