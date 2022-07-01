package share

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MustConnectDb() *gorm.DB {
	uname := MustGetEnv("MYSQL_USER")
	pass := MustGetEnv("MYSQL_PASSWORD")
	dbHost := MustGetEnv("DB_HOST")
	dbName := MustGetEnv("DB_NAME")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", uname, pass, dbHost, dbName)

	logger := createLogger()

	gormDb, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: logger})

	if err != nil {
		log.Print(err)
		ExitFatal("error connecting to the database")
	}
	db, err := gormDb.DB()
	if err != nil {
		log.Print(err)
		ExitFatal("error connecting to the database")
	}
	db.Ping()
	if err != nil {
		log.Print(err)
		ExitFatal("error connecting to the database")
	}
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Hour)

	return gormDb
}

func createLogger() logger.Interface {

	logger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})
	return logger
}
