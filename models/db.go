package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"fmt"
)

func InitDb(host, user, password, dbName string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s sslmode=disable user=%s password=%s dbname=%s", host, user, password, dbName,
	)

	db, err := gorm.Open("postgres", connectionString)

	db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func InitTestDb(host, user, password, dbName string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s sslmode=disable user=%s password=%s dbname=%s_test", host, user, password, dbName,
	)

	db, err := gorm.Open("postgres", connectionString)

	db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func CreateTables(db *gorm.DB) {
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}

	if !db.HasTable(&TodoList{}) {
		db.CreateTable(&TodoList{})
		db.Model(&Todo{}).AddForeignKey(
			"user_id", "users(id)", "CASCADE", "RESTRICT",
		)
	}

	if !db.HasTable(&Todo{}) {
		db.CreateTable(&Todo{})
		db.Model(&Todo{}).AddForeignKey(
			"user_id", "users(id)", "CASCADE", "RESTRICT",
		)
		db.Model(&Todo{}).AddForeignKey(
			"todo_list_id", "todo_lists(id)", "CASCADE", "RESTRICT",
		)
	}

	db.AutoMigrate(&User{}, &TodoList{}, &Todo{})
}

func DropTables(db *gorm.DB) {
	db.Delete(&User{})
	db.Delete(&TodoList{})
	db.Delete(&Todo{})
}
