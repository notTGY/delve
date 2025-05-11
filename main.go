package main

import (
	"fmt"
	"log"
	"os"
  "time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"

  "github.com/nottgy/delve/llm"
  "github.com/nottgy/delve/memory"
)

func main() {
	_ = godotenv.Load()
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10*time.Second},
	})
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userMessage := c.Text()
    user := c.Sender()
    s := fmt.Sprintf("%s: %s\n", user.Username, userMessage)

    dialog, err := memory.Query(user.ID)
		if err != nil {
			log.Println("Memory error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}

    messages := append(
      dialogToMessages(dialog),
      userMessage,
    )

    embStart := time.Now()
		_, err = llm.Embed(messages)
		if err != nil {
			log.Println("OpenAI error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}
    s += fmt.Sprintf("Embed %d in %s\n", len(messages), time.Since(embStart))

    genStart := time.Now()
		resp, err := llm.Message(userMessage)
		if err != nil {
			log.Println("OpenAI error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}
    s += fmt.Sprintf("Gen in %s\n", time.Since(genStart))

    err = memory.Save(user.ID, userMessage, resp)
		if err != nil {
			log.Println("Memory error:", err)
		}

    fmt.Printf("%s------------------\n", s)
		return c.Send(resp)
	})

	fmt.Println("Bot is running...")
	bot.Start()
}

func dialogToMessages(dialog []memory.Dialog) []string {
  messages := []string{}
  for _, d := range dialog {
    messages = append(
      messages,
      d.UserMessage,
      d.AssistantMessage,
    )
  }
  return messages
}
