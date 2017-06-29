package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"fmt"
)

func InitDb(user, password, dbname string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=localhost sslmode=disable user=%s password=%s dbname=%s", user, password, dbname,
	)

	db, err := gorm.Open("postgres", connectionString)

	db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func InitTestDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./test_data.db")
	db.LogMode(false)

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

	if !db.HasTable(&Todo{}) {
		db.CreateTable(&Todo{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Todo{})
	}
}

func DropTables(db *gorm.DB) {
	db.DropTableIfExists(&TodoList{})
	db.DropTableIfExists(&Todo{})
}
