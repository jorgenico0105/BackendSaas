package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas-medico/internal/config"
)

var DB *gorm.DB

func Connect() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var gormConfig *gorm.Config
	if cfg.Environment == "development" {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	} else {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established successfully")
}

func GetDB() *gorm.DB {
	return DB
}
