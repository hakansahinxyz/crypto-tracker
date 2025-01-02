package services

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramService(botToken string, chatID int64) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &TelegramService{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramService) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
	return err
}
