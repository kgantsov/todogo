package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"fmt"
)

func InitDb(host, user, password, dbName string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s sslmode=disable user=%s password=%s dbName=%s", host, user, password, dbName,
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
	}

	if !db.HasTable(&Todo{}) {
		db.CreateTable(&Todo{})
		db.Model(&Todo{}).AddForeignKey(
			"todo_list_id", "todo_lists(id)", "CASCADE", "RESTRICT",
		)
	}
}

func DropTables(db *gorm.DB) {
	db.DropTableIfExists(&TodoList{})
	db.DropTableIfExists(&Todo{})
}
