package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var RequestError = errors.New("ошибка при попытке выполнить запрос к серверу")

func UserIsBot(userID string) (bool, error) {
	url := fmt.Sprintf("https://api.telegram.org/%s/getChat?chat_id=%s", BuildToken(), userID)

	resp, err := http.Get(url)
	if err != nil {
		return false, RequestError
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil
	}

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			ID       int64  `json:"id"`
			IsBot    bool   `json:"is_bot"`
			Username string `json:"username"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if result.Ok {
		if result.Result.IsBot {
			return true, nil
		} else {
			return false, nil
		}
	}
	return true, fmt.Errorf("failed to get chat info")

}
