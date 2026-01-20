package config

import (
	"fmt"
	"log"

	"2026-FM247-BackEnd/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	log.Println("Database connection established")

	// 自动迁移
	db.AutoMigrate(
		&models.User{},
		&models.TotalStudyData{},
		&models.DailyStudyData{},
		&models.MonthlyStudyData{},
		&models.Todo{},
		&models.Note{},
		&models.TokenBlacklist{},
	)
	log.Println("Database migrated successfully")
	return nil
}
