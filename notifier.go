package main

import (
	"fmt"
	"net/http"
)

type TelegramNotifier struct {
	token  string
	chatID string
}

func NewTelegramNotifier(
	token string,
	chatID string,
) *TelegramNotifier {
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
