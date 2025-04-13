package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/repository/database"
	"github.com/spf13/cobra"
)

var migrateDataCmd = &cobra.Command{
	Use:   "migrate",
	Short: "init database",
	Long:  "init database",
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.NewConfig()
		dbConfig := conf.Database
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
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Printf("error connecting to the database: %+v\n", err)
			return
		}
		dbMigrator, err := database.NewMigrator(db)
		if err != nil {
			return
		}
		err = dbMigrator.Up(context.Background())
		if err != nil {
			return
		}

		log.Println("Migrate data successfully")
	},
}
