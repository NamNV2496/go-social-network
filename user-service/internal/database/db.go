package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // Import MySQL goqu dialect
	_ "github.com/go-sql-driver/mysql"               // Import MySQL driver
	"github.com/namnv2496/user-service/internal/configs"
)

type Database struct {
	*sql.DB
}

func NewDatabase(
	dbConfig configs.Database,
) (*sql.DB, func(), error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	fmt.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("error connecting to the database: %+v\n", err)
		return nil, nil, err
	}

	dbMigrator, err := NewMigrator(db)
	if err != nil {
		return db, nil, err
	}
	err = dbMigrator.Up(context.Background())
	if err != nil {
		return db, nil, err
	}
	cleanup := func() {
		db.Close()
	}
	return db, cleanup, nil
}

func InitializeGoquDB(db *sql.DB) *goqu.Database {
	return goqu.New("mysql", db)
}
