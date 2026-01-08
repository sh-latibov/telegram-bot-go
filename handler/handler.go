package handler

import (
	"fmt"
	"log"
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sh-latibov/telegram-bot-go/clients/openweather"
)

type Handler struct {
	bot        *tgbotapi.BotAPI
	owClient   *openweather.OpenWeatherClient
	userCities map[int64]string
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:        bot,
		owClient:   owClient,
		userCities: make(map[int64]string),
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil { // If we got a message
		return
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "city":
			h.handleCityCommand(update)
			return
		case "weather":
			h.handleWeatherCommand(update)
			return
		default:
			h.handleUnknownCommand(update)
			return
		}

	}

}

func (h *Handler) Start() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}

func (h *Handler) handleWeatherCommand(update tgbotapi.Update) {
	city, ok := h.userCities[update.Message.From.ID]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала установите город командой /city Ташкент")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	coordinates, err := h.owClient.Coordinates(city)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли найти координаты")
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Printf("[ERROR] Не удалось получить погоду для координат (lat: %.2f, lon: %.2f): %v\n", coordinates.Lat, coordinates.Lon, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить данные о погоде. Попробуйте позже.")
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Температура в %s: %d°C", city, int(math.Round(weather.Temp))))
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)

}

func (h *Handler) handleCityCommand(update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	h.userCities[update.Message.Chat.ID] = city
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город %s сохранен", city))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleUnknownCommand(update tgbotapi.Update) {
	log.Printf("[WARNING] Неизвестная команда от пользователя @%s (ID: %d): '%s'\n", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такая команда не поддерживается. Используйте /city или /weather")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}
