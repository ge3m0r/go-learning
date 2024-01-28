package dao

import (
	"basic-go/webook/internal/domain"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{})
}
