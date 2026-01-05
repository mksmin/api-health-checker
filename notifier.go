package main

import (
	"fmt"
	"net/http"
	"os"
)

type TelegramNotifier struct {
	token  string
	chatID string
}

func NewTelegramNotifier() *TelegramNotifier {
	token := os.Getenv("TG_BOT_TOKEN")
	chatID := os.Getenv("TG_CHAT_ID")

	if token == "" || chatID == "" {
		panic("TG_BOT_TOKEN and TG_CHAT_ID must be set")
	}

	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
	}
}

func (
	t *TelegramNotifier,
) sendMessage(
	text string,
) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		t.token, t.chatID, text,
	)
	_, _ = http.Get(url)
}

func (
	t *TelegramNotifier,
) NotifyDown(
	s *Service,
) {
	t.sendMessage(
		fmt.Sprintf(
			"⚠️ Service %s is down",
			s.Name,
		),
	)
	LogEvent(fmt.Sprintf("Send message about %s down", s.Name))

}
func (
	t *TelegramNotifier,
) NotifyUp(
	s *Service,
) {
	t.sendMessage(
		fmt.Sprintf(
			"✅️ Service %s is recover",
			s.Name,
		),
	)
	LogEvent(fmt.Sprintf("Send message about %s recover", s.Name))

}
