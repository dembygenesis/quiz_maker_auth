package main

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/cmd/database/migrations/migrate"
	"github.com/dembygenesis/quiz_maker_auth/cmd/database/migrations/seed"
	db2 "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/db"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	fmt.Println("============Loaded============")

	if err != nil {
		panic("Error loading .env files:" + err.Error())
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var err error
	db, err := db2.GetGormInstance(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)

	if err != nil {
		return errors.New("error fetching gorm instance: " + err.Error())
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = migrate.Run(tx)
		if err != nil {
			return err
		}
		err = seed.Run(tx)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}