package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SendText sends one plain text message without retries.
func SendText(ctx context.Context, token, chatID, text string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if strings.TrimSpace(chatID) == "" {
		return errors.New("TELEGRAM_CHAT_ID is required")
	}

	form := url.Values{"chat_id": {chatID}, "text": {text}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.telegram.org/bot"+token+"/sendMessage", strings.NewReader(form.Encode()))
	if err != nil {
		// Request errors can contain the URL, which includes the secret token.
		return errors.New("could not construct Telegram request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		// Do not expose the underlying error: it can include the token-bearing URL.
		return errors.New("Telegram request failed (network error, timeout, or cancellation)")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Telegram HTTP request failed with status %d", res.StatusCode)
	}

	var response struct {
		OK        bool `json:"ok"`
		ErrorCode int  `json:"error_code"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return errors.New("could not decode Telegram response")
	}
	if !response.OK {
		// Report only the code, avoiding untrusted response text containing secrets.
		return fmt.Errorf("Telegram API rejected the message (error code %d)", response.ErrorCode)
	}
	return nil
}
