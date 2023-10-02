package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	config "github.com/mochammadshenna/arch-pba-template/config"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

func TestOpenConnection(t *testing.T) {
	var dbConfig = config.Get().Database

	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
	)

	db, err := sql.Open("mysql", mysqlInfo)
	panicOnError(err)
	if err = db.Ping(); err != nil {
		logger.Fatal(context.TODO(), err)
		if err = db.Close(); err != nil {
			logger.Fatal(context.TODO(), err)
		}
	}

	defer db.Close()
}

func panicOnError(err error) {
	if err != nil {
		log.Printf("panic on config %v", err)
		panic(err)
	}
}
