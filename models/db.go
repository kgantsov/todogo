package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./data.db")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	return db
}

func CreateTables(db *gorm.DB) {
	if !db.HasTable(&TodoList{}) {
		db.CreateTable(&TodoList{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&TodoList{})
	}
}

func DropTables(db *gorm.DB) {
	db.DropTableIfExists(&TodoList{})
	db.DropTableIfExists(&Todo{})
}
