package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
)

func main() {
	_ = godotenv.Load()
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	config := openai.DefaultConfig("ollama")
	config.BaseURL = "http://localhost:11434/v1/"
	client := openai.NewClientWithConfig(config)

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userMessage := c.Text()

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
			log.Println("OpenAI error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}

		response := resp.Choices[0].Message.Content
		_, answer, is_thinking := strings.Cut(response, "</think>")
		if is_thinking {
			return c.Send(strings.TrimSpace(answer))
		}
		return c.Send(response)
	})

	fmt.Println("Bot is running...")
	bot.Start()
}
