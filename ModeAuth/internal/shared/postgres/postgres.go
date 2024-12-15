package postgres

import (
	"ModeAuth/internal/shared/config"
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() (*gorm.DB, error) {
	log.Println(logging.INFO + "Connection to Postgres")

	dsn, err := config.GetPostgres()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn.ConnStr), &gorm.Config{})
	if err != nil {
		log.Printf(logging.ERROR+"Failed to connect to Postgres: %v", err)

		return nil, err
	}

	log.Printf(logging.INFO + "Connection to Postgres successful")

	if err := db.AutoMigrate(&dto.User{}, &dto.BlockedUser{}, &dto.Report{}); err != nil {
		return nil, err
	}

	log.Println(logging.INFO + "Tables User, BlockUser, Report - loaded successful into Postgres")

	return db, nil
}
