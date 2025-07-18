package infrastructure

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/steradian_interview?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}
