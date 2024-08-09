package extension

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataBase struct {
	db *gorm.DB
}

func (m *DataBase) Init() error {
	if m.db != nil {
		return nil
	}
	env := os.Getenv("ENV")
	database := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	maxIdleConns := 2
	maxOpenConns := 100
	maxLifeTime := 3600
	if len(os.Getenv("DB_MAX_IDLE_CONNS")) > 0 {
		maxIdleConns, _ = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
	}
	if len(os.Getenv("DB_MAX_OPEN_CONNS")) > 0 {
		maxOpenConns, _ = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	}
	if len(os.Getenv("DB_MAX_LIFE_TIME")) > 0 {
		maxLifeTime, _ = strconv.Atoi(os.Getenv("DB_MAX_LIFE_TIME"))
	}

	dsnAppend := os.Getenv("DB_DSN_APPEND")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", username, password, host, port, database, dsnAppend)
	mysqlDb := mysql.Open(dsn)
	ormConfig := gorm.Config{}
	if env == "production" {
		ormConfig.Logger = logger.Default.LogMode(logger.Error)
	} else {
		ormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	linkDb, err := gorm.Open(mysqlDb, &ormConfig)
	if err != nil {
		log.Println(err)
		return err
	}
	sqlDB, err := linkDb.DB()
	if err != nil {
		log.Println(err)
		return err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns) // idle connect num
	sqlDB.SetMaxOpenConns(maxOpenConns) // max connect num
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)
	m.db = linkDb

	return m.testConnect()
}

func (m *DataBase) Db() *gorm.DB {
	return m.db
}

func (m *DataBase) testConnect() error {
	var sum int
	m.db.Raw("SELECT 1+1").Scan(&sum)
	log.Println("db test: 1+1 =", sum)
	if sum != 2 {
		return errors.New("test fail")
	}
	return nil
}
