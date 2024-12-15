package utils

import (
	"ModeAuth/pkg/logging"
	"flag"
	"fmt"
	"log"
)

var BotToken = flag.String("token", "", "Telegram bot token")

func GetToken() {
	log.Println(logging.INFO + "Starting getting the Telegram token from the flag")

	flag.Parse()

	if *BotToken == "" {
		fmt.Println("The bot token was not passed in the launch parameters. Please provide the bot token using the \" -token \" flag")
		return
	}

	log.Println(logging.INFO + "Telegram token successfully loaded")

	CheckBotToken(*BotToken)
}

func BuildToken() string {
	return fmt.Sprintf("bot%s", *BotToken)
}
