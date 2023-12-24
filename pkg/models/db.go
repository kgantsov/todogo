package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(host, user, password, dbName string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s sslmode=disable user=%s password=%s dbname=%s", host, user, password, dbName,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	// db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func InitDbURI(connectionString string, debug bool) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func InitTestDb(host, user, password, dbName string, debug bool) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s sslmode=disable user=%s password=%s dbname=%s_test", host, user, password, dbName,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	// db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func InitTestDbURI(connectionString string, debug bool) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	// db.LogMode(debug)

	if err != nil {
		panic(err)
	}

	return db
}

func CreateTables(db *gorm.DB) {
	if !db.Migrator().HasTable(&User{}) {
		db.Migrator().CreateTable(&User{})
	}

	if !db.Migrator().HasTable(&TodoList{}) {
		db.Migrator().CreateTable(&TodoList{})
		// db.Model(&TodoList{}).AddForeignKey(
		// 	"user_id", "users(id)", "CASCADE", "RESTRICT",
		// )
	}

	if !db.Migrator().HasTable(&Todo{}) {
		db.Migrator().CreateTable(&Todo{})
		// db.Model(&Todo{}).AddForeignKey(
		// 	"user_id", "users(id)", "CASCADE", "RESTRICT",
		// )
		// db.Model(&Todo{}).AddForeignKey(
		// 	"todo_list_id", "todo_lists(id)", "CASCADE", "RESTRICT",
		// )
	}

	db.AutoMigrate(&User{}, &TodoList{}, &Todo{})
}

func DropTables(db *gorm.DB) {
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&TodoList{})
	db.Migrator().DropTable(&Todo{})
}
