package database

import (
	"book-search/webcrawler/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBClient() (*gorm.DB, error) {
	dsn := "host=localhost user=admin password=1q2w3e4r dbname=book_search port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = migrateTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateTables(db *gorm.DB) error {
	return db.AutoMigrate(&models.Book{}, &models.Author{}, &models.AuthorBook{})
}

func CloseDBClient(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
