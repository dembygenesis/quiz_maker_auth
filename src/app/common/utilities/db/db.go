package db

import (
	UtilitiesString "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/database"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
)

func GetArrayOfIntAsWhereIn(i []int) string {
	sqlCondition := "("

	for _, v := range i {
		sqlCondition += `"` + strconv.Itoa(v) + `",`
	}


	sqlCondition = UtilitiesString.TrimSuffix(sqlCondition, ",")

	sqlCondition += ")"

	return sqlCondition
}

func ValidEntry(client string, v string, c string, t string) (bool, error) {
	var count int
	var hasEntry bool

	sql := `
		SELECT COUNT(*) FROM ` + t + `
		WHERE 1 = 1
			AND ` + c + ` = ?
	`

	err := database.DBInstancePublic.Get(&count, sql, v)

	if err != nil {
		return hasEntry, err
	}

	if count > 0 {
		hasEntry = true
	} else {
		hasEntry = false
	}

	return hasEntry, err
}

func GetLastInsertID() (int, error) {
	var id int

	sql := "SELECT LAST_INSERT_ID()"

	err := database.DBInstancePublic.Get(&id, sql)

	return id, err
}

func GetLastInsertIDTx(t *sqlx.Tx) (int, error) {
	var id int

	sql := "SELECT LAST_INSERT_ID()"

	err := t.Get(&id, sql)

	return id, err
}

// GetGormInstance - returns a gorm instance
func GetGormInstance(
	dbHost string,
	dbUser string,
	dbPassword string,
	dbDatabase string,
	dbPort string,
) (*gorm.DB, error) {
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			// SingularTable: true,
		},
	})
	return db, err
}

func GetLastInsertIDGorm(tx *gorm.DB) (int, error) {
	var lastInsertId int
	err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&lastInsertId).Error
	return lastInsertId, err
}