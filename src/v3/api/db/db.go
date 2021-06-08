package db

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"strings"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var (
	// Handle to DB connection
	Handle *gorm.DB
)

func init() {
	fmt.Println("============================== RUNNING DB HANDLER ==============================")

	var dbConfig config.DB
	err := envconfig.Process("partner", &dbConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=%s&parseTime=True&loc=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Name,
		dbConfig.Charset,
		strings.Replace(dbConfig.Timezone, "/", "%2F", -1))

	// Let's see you implement
	/*namingStrategy := schema.NamingStrategy{
		SingularTable: false,
	}*/

	Handle, err = gorm.Open(mysql.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// NamingStrategy: namingStrategy,
	})

	if err != nil {
		log.Fatal("Could not connect database")
	} else {
		log.Println("============== Successfully connected ==============", Handle)
	}

	Handle.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
}
