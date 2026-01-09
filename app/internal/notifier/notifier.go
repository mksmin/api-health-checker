package notifier

import (
	"fmt"
	"healthchecker/internal/common"
	"healthchecker/internal/logs"
	"net/http"
	"net/url"
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
	safeText := url.QueryEscape(text)
	urlTg := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		t.token, t.chatID, safeText,
	)
	resp, err := http.Get(urlTg)
	if err != nil {
		logs.LogEvent(
			fmt.Sprintf(
				"Error while sending a message: %v",
				err,
			),
		)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logs.LogEvent("Telegram returned non-200 status code")
	}
}

func (
	t *TelegramNotifier,
) NotifyDown(
	s *common.Service,
) {
	t.sendMessage(
		fmt.Sprintf(
			"⚠️ Service %s is down",
			s.Name,
		),
	)
	logs.LogEvent(fmt.Sprintf("Send message about %s down", s.Name))

}

func (
	t *TelegramNotifier,
) NotifyUp(
	s *common.Service,
) {
	t.sendMessage(
		fmt.Sprintf(
			"✅️ Service %s is recover",
			s.Name,
		),
	)
	logs.LogEvent(fmt.Sprintf("Send message about %s recover", s.Name))

}
