package main

import (
	"JsonParse/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	Controller api.Controller
)

func init() {
	InitDB()
	InitCon()
	InitViper()
}

func main() {
	fmt.Println(viper.GetString("itt.url.class"))
	r := gin.New()
	an := r.Group("/answer")
	{
		an.GET("/init", Controller.InitTable)
		an.GET("/data", Controller.GetAnswerByClass)
		an.GET("/question", Controller.GetQuestionsByClass)
		an.GET("/abq", Controller.GetAnswerByClassAndQuestionId)
		an.POST("/upload", Controller.UploadFile)
		an.POST("/sync", Controller.AutoSyncAnswer)
	}

	cl := r.Group("/class")
	{
		cl.GET("/getList", Controller.GetClassList)
		cl.GET("/addClass", Controller.AddClass)
		cl.GET("/delete", Controller.DeleteClassByCode)
	}

	_ = r.Run()
	return
}
