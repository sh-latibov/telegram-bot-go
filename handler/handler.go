package handler

import (
	"fmt"
	"log"
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sh-latibov/telegram-bot-go/clients/openweather"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}
}

func (h *Handler) HandUpdate(update tgbotapi.Update) {
	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		coordinates, err := h.owClient.Coordinates(update.Message.Text)
		if err != nil {
			log.Printf("weather error: %+v\n", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли найти координаты")
			msg.ReplyToMessageID = update.Message.MessageID

			h.bot.Send(msg)
			return
		}

		weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли найти погоду")
			msg.ReplyToMessageID = update.Message.MessageID

			h.bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Температура в %s: %d°C", update.Message.Text, int(math.Round(weather.Temp))))
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)
	}
}
