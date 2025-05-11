package llm

import (
	"github.com/sashabaranov/go-openai"
  "context"
  "strings"
)

var config = openai.DefaultConfig("ollama")
var client *openai.Client

func init() {
  config.BaseURL = "http://localhost:11434/v1/"
  client = openai.NewClientWithConfig(config)
}

func Embed(inputs []string) ([]openai.Embedding, error) {
  resp, err := client.CreateEmbeddings(
    context.Background(),
    openai.EmbeddingRequest{
      Input: inputs,
      Model: "snowflake-arctic-embed:22m",
    },
  )
  if err != nil {
    return nil, err
  }
  return resp.Data, nil
  /*
  queryEmbedding := queryResponse.Data[0]
	targetEmbedding := targetResponse.Data[0]

	similarity, err := queryEmbedding.DotProduct(&targetEmbedding)
	if err != nil {
		log.Fatal("Error calculating dot product:", err)
	}

	log.Printf("The similarity score between the query and the target is %f", similarity)
  */
}

func Message(userMessage string) (string, error) {
  resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
      Model: "qwen3:0.6b",
      Messages: []openai.ChatCompletionMessage{
        {
          Role:    openai.ChatMessageRoleSystem,
          Content: "You are a helpful AI assistant, called Jenny Tull. You were trained by company @gamma_code",
        },
        {
          Role:    openai.ChatMessageRoleUser,
          Content: userMessage,
        },
      },
    },
  )
  if err != nil {
    return "", err
  }

  response := resp.Choices[0].Message.Content
  _, answer, is_thinking := strings.Cut(response, "</think>")
  if is_thinking {
    return strings.TrimSpace(answer), nil
  }
  return response, nil
}
