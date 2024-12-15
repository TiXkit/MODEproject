package utils

import (
	"ModeAuth/pkg/logging"
	"fmt"
	"log"
	"net/http"
)

func CheckBotToken(token string) {
	log.Println(logging.INFO + "Starting process of token validity checks")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token), nil)
	if err != nil {
		log.Fatal(logging.ERROR+"Failed to create request to verify bot token: ", err)

	}

	log.Println(logging.INFO + "Bot token verification request header set")
	req.Header.Set("Authorization", "Bot "+token)

	log.Println(logging.INFO + "A request is made to verify the bot token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(logging.ERROR+"Failed to make a request to the server to check the validity of the bot token: ", err)

	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(logging.WARN + "The token transferred when the application is launched is not valid")

	}

	log.Println(logging.INFO + "The token transferred when the application is launched is valid")
}
