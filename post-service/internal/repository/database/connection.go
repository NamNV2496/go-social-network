package database

import (
	"fmt"
	"log"
	"os"

	"github.com/namnv2496/post-service/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabaseConnection(
	dbConfig configs.Database,
) *gorm.DB {
	var connectionString string
	if value := os.Getenv("DATABASE_URL"); value != "" {
		connectionString = value
	} else {
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
		)
	}

	fmt.Println(connectionString)
	db, err := gorm.Open(
		mysql.Open(connectionString),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: false,
			PrepareStmt:                              false,
		},
	)
	if err != nil {
		log.Printf("error connecting to the database: %+v\n", err)
		return nil
	}

	return db
}
