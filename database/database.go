package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	config "github.com/mochammadshenna/arch-pba-template/config"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

func NewDB() *sql.DB {
	return newDb(config.Get().Database.DbName)
}

func newDb(dbName string) *sql.DB {
	var dbConfig = config.Get().Database

	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbName,
	)

	db, err := sql.Open("mysql", mysqlInfo)
	panicOnError(err)
	if err = db.Ping(); err != nil {
		logger.Fatal(context.TODO(), err)
		if err = db.Close(); err != nil {
			logger.Fatal(context.TODO(), err)
		}
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func panicOnError(err error) {
	if err != nil {
		log.Printf("panic on config %v", err)
		panic(err)
	}
}
