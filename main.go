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

    memory.Query(struct{ID int64}{user.ID})

    embStart := time.Now()
		_, err := llm.Embed([]string{userMessage})
		if err != nil {
			log.Println("OpenAI error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}
    s += fmt.Sprintf("Embed in %s\n", time.Since(embStart))

    genStart := time.Now()
		resp, err := llm.Message(userMessage)
		if err != nil {
			log.Println("OpenAI error:", err)
			return c.Send("Sorry, I couldn't process your request.")
		}
    s += fmt.Sprintf("Gen in %s\n", time.Since(genStart))
    fmt.Printf("%s------------------\n", s)
		return c.Send(resp)
	})

	fmt.Println("Bot is running...")
	bot.Start()
}
