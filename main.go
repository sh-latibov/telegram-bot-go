package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sh-latibov/telegram-bot-go/clients/openweather"
	"github.com/sh-latibov/telegram-bot-go/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	owClient := openweather.New(os.Getenv("WEATHER_KEY"))

	botHandlaer := handler.New(bot, owClient)

	botHandlaer.Start()

}
