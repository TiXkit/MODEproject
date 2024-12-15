package config

import (
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"log"
	"os"
)

func GetPostgres() (*dto.PostgresConfig, error) {
	log.Println(logging.INFO + "Starting the process of obtaining the config for Postgres")

	connStrPg := os.Getenv("POSTGRES_CONNECT_STRING")
	if connStrPg == "" {
		log.Fatal(logging.FATAL + "Failed to get Postgres connection information from .env file")
	}

	log.Println(logging.INFO + "Config for Postgres successfully received")

	return &dto.PostgresConfig{
		ConnStr: connStrPg,
	}, nil
}

func GetRedis() *dto.RedisConfig {
	log.Println(logging.INFO + "Starting the process of obtaining the config for Redis")

	addrRedis := os.Getenv("REDIS_ADDR")
	passRedis := os.Getenv("REDIS_PASSWORD")
	if addrRedis == "" || passRedis == "" {
		log.Printf("Failed to get Redis connection information from .env file")
	}

	log.Println(logging.INFO + "Config for Redis successfully received")

	return &dto.RedisConfig{
		Addr:     addrRedis,
		Password: passRedis,
	}
}

func GetAddress() *dto.AddressHTTP {
	log.Println(logging.INFO + "Starting the process of obtaining the config for HTTP")

	addrHTTP := os.Getenv("HTTP_ADDRESS")
	if addrHTTP == "" {
		addrHTTP = ":8086"
	}

	log.Println(logging.INFO + "Config for HTTP successfully received")

	return &dto.AddressHTTP{
		Addr: addrHTTP,
	}
}
