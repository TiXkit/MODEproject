package main

import (
	"ModeAuth/internal/repository"
	"ModeAuth/internal/service"
	"ModeAuth/internal/shared/config"
	"ModeAuth/internal/shared/postgres"
	"ModeAuth/internal/shared/redis"
	"ModeAuth/internal/shared/utils"
	"ModeAuth/internal/transport"
	"ModeAuth/pkg/logging"
	"context"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(logging.FATAL + "No .env file found")
	}
	log.Println(logging.INFO + ".env file loaded")
}

func main() {
	utils.GetToken()

	log.Println(logging.INFO + "New context created")
	ctx := context.Background()

	log.Println(logging.INFO + "Context timeout set to five")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal(logging.FATAL+"Error initializing database: ", err)
	}

	rDB := redis.ConnectRedis(ctx)

	repo := repository.NewRepository(db, rDB)

	auth := service.NewAuthentication(repo)
	state := service.NewState(repo)

	addrHTTP := config.GetAddress()

	transport.RunRouter(auth, state, addrHTTP.Addr)
}
