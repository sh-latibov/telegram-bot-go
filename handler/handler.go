package handler

import (
	"context"
	"fmt"
	"log"
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sh-latibov/telegram-bot-go/clients/openweather"
	"github.com/sh-latibov/telegram-bot-go/models"
)

type userRepository interface {
	SaveUser(ctx context.Context, userID int64) error
	SaveUserCity(ctx context.Context, userID int64, city string) error
	GetUserCity(ctx context.Context, userID int64) (string, error)
	UpdateUserCity(ctx context.Context, userID int64, city string) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	IsUserExists(ctx context.Context, userID int64) (bool, error)
}

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	userRepo userRepository
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, userRepo userRepository) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		userRepo: userRepo,
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil { // If we got a message
		return
	}

	ctx := context.Background()

	if update.Message.IsCommand() {
		err := h.ensureUser(ctx, update)
		if err != nil {
			log.Printf("[ERROR] Не удалось обеспечить наличие пользователя в базе для @%s (ID: %d): %v\n", update.Message.From.UserName, update.Message.From.ID, err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обработке вашей команды. Пожалуйста, попробуйте еще раз.")
			msg.ReplyToMessageID = update.Message.MessageID

			h.bot.Send(msg)
			return
		}

		switch update.Message.Command() {
		case "city":
			h.handleCityCommand(ctx, update)
			return
		case "weather":
			h.handleWeatherCommand(ctx, update)
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

func (h *Handler) handleWeatherCommand(ctx context.Context, update tgbotapi.Update) {
	city, err := h.userRepo.GetUserCity(ctx, update.Message.Chat.ID)
	if err != nil {
		log.Printf("[ERROR] Не удалось получить город для пользователя @%s (ID: %d): %v\n", update.Message.From.UserName, update.Message.From.ID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить ваш город. Пожалуйста, установите его с помощью команды /city <название_города>.")
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)
		return
	}

	if city == "" {
		log.Printf("[ERROR] Не удалось получить город для пользователя @%s (ID: %d): %v\n", update.Message.From.UserName, update.Message.From.ID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить ваш город. Пожалуйста, установите его с помощью команды /city <название_города>.")
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

func (h *Handler) handleCityCommand(ctx context.Context, update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	err := h.userRepo.UpdateUserCity(ctx, update.Message.From.ID, city)
	if err != nil {
		log.Printf("[ERROR] Не удалось сохранить город для пользователя @%s (ID: %d): %v\n", update.Message.From.UserName, update.Message.From.ID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сохранить город. Попробуйте еще раз.")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
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

func (h *Handler) ensureUser(ctx context.Context, update tgbotapi.Update) error {
	user, err := h.userRepo.GetUser(ctx, update.Message.From.ID)
	if err != nil {

		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		err := h.userRepo.SaveUser(ctx, update.Message.From.ID)
		if err != nil {

			return fmt.Errorf("error saving user: %w", err)
		}
	}
	return nil
}
