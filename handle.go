package main

import (
	"JsonParse/api"
	Config "JsonParse/config"
	"JsonParse/model"
	"JsonParse/repository"
	"JsonParse/service"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {
	DB = model.GetSqliteDB()
}

func InitViper() {
	if err := Config.Init("./config/config.yaml"); err != nil {
		panic(err)
	}
}

func InitCon() {
	Controller = api.Controller{
		AnswerService: &service.AnswerService{
			AnswerRepo: &repository.AnswerRepo{
				DB: model.DataBase{
					SqlLite: DB,
				},
			},
		},
	}
}
