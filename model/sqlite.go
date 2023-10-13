package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type DataBase struct {
	SqlLite *gorm.DB
}

var DB *DataBase

func (db *DataBase) Init() {
	DB = &DataBase{
		SqlLite: GetSqliteDB(),
	}
}

func GetSqliteDB() *gorm.DB {
	return InitSelfDB()
}

func InitSelfDB() *gorm.DB {
	db := openDB("test.db")
	return db
}

func openDB(database string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("Database connection failed. Database name: %s,Error:%s", database, err.Error())
	}
	_ = db.AutoMigrate(Answer{}, Class{})
	setupDB(db)
	return db
}

func setupDB(db *gorm.DB) {
	// 用于设置闲置的连接数.
	sql, _ := db.DB()
	sql.SetMaxIdleConns(5)
}
