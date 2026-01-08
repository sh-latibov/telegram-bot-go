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
		log.Printf("[WARNING] Ошибка загрузки файла .env: %v. Попытаюсь использовать переменные окружения системы.\n", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("[ERROR] Ошибка инициализации бота: %v. Убедитесь, что переменная BOT_TOKEN установлена и корректна.\n", err)
	}

	bot.Debug = true

	log.Printf("[INFO] Бот успешно авторизован. Аккаунт: @%s (ID: %d)\n", bot.Self.UserName, bot.Self.ID)

	weatherKey := os.Getenv("WEATHER_KEY")
	if weatherKey == "" {
		log.Fatal("[ERROR] Переменная окружения WEATHER_KEY не установлена")
	}
	owClient := openweather.New(weatherKey)
	log.Println("[INFO] OpenWeather клиент инициализирован")

	botHandlaer := handler.New(bot, owClient)
	log.Println("[INFO] Обработчик событий инициализирован")

	log.Println("[INFO] Бот начал слушать входящие сообщения...")
	botHandlaer.Start()

}
