package llm

import (
	"context"
	"slices"
	"sort"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/nottgy/delve/memory"
)

var config = openai.DefaultConfig("ollama")
var client *openai.Client

const system_prompt = `You are a helpful AI assistant, called Jenny Tull.
You were trained by company @gamma_code`

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
}

func Message(userMessage string, relevantDialogIdxs []int, dialog []memory.Dialog) (string, error) {
	messages := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: system_prompt,
	}}

	for _, idx := range relevantDialogIdxs {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: dialog[idx].UserMessage,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: dialog[idx].AssistantMessage,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "qwen3:0.6b",
			Messages: messages,
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

func TopK(embs []openai.Embedding, target openai.Embedding, k int) ([]int, error) {
	type pair struct {
		index int
		score float32
	}

	pairs := make([]pair, len(embs))
	for i, emb := range embs {
		score, err := emb.DotProduct(&target)
		if err != nil {
			return nil, err
		}
		size, err := emb.DotProduct(&emb)
		if err != nil {
			return nil, err
		}
		pairs[i] = pair{index: i / 2, score: score / size}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	n := len(pairs)
	if k > n/2 {
		k = n / 2
	}

	result := make([]int, k)
	for i, j := 0, 0; i < k; i++ {
		for ; j < n; j++ {
			idx := pairs[j].index
			if slices.Index(result, idx) == -1 {
				result[i] = idx
				break
			}
		}
	}

	slices.Sort(result)

	return result, nil
}
