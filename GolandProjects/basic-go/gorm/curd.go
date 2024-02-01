package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	Code  string
	Price uint
	gorm.Model
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db = db.Debug()

	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "D42", Price: 100})

}
