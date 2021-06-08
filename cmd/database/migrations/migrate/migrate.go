package migrate

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/modelsV2"
	"gorm.io/gorm"
	"os"
)

func Run(tx *gorm.DB) error {
	fmt.Println("----------- Migrate -----------")
	var err error
	err = tx.Exec("DROP DATABASE " + os.Getenv("DB_DATABASE") + ";").Error
	if err != nil {
		return err
	}
	err = tx.Exec("CREATE DATABASE " + os.Getenv("DB_DATABASE") + ";").Error
	if err != nil {
		return err
	}
	err = tx.Exec("USE " + os.Getenv("DB_DATABASE") + ";").Error
	if err != nil {
		return err
	}
	err =  tx.AutoMigrate(
		&modelsV2.UserType{},
		&modelsV2.User{},
		&modelsV2.Quiz{},
		&modelsV2.QuizQuestion{},
		&modelsV2.QuizChoice{},
	)
	if err != nil {
		return err
	}
	return nil
}
