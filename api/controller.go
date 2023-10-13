package api

import (
	"JsonParse/model"
	"JsonParse/service"
	"JsonParse/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Controller struct {
	AnswerService *service.AnswerService
}

// 初始化数据
func (con *Controller) InitTable(context *gin.Context) {
	con.AnswerService.InitAnswer()
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   "",
			"code":   http.StatusOK,
		},
	)
}

// 根据科目获取答案
func (con *Controller) GetAnswerByClass(context *gin.Context) {
	Class := context.Query("class")
	answerList := con.AnswerService.GetAnswerByClass(Class)
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   answerList,
			"code":   http.StatusOK,
		},
	)
}

// 上传文件
func (con *Controller) UploadFile(context *gin.Context) {
	//form表单
	file, header, err := context.Request.FormFile("upload")
	if err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": fmt.Sprintf("上传文件失败: %s", err.Error()),
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
		return
	}
	defer file.Close()
	err = con.AnswerService.JsonParse(&file, header.Filename)

	if err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": fmt.Sprintf("上传文件失败: %s", err.Error()),
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   "",
			"code":   http.StatusOK,
		},
	)
	return
}

// 获取问题列表
func (con *Controller) GetQuestionsByClass(context *gin.Context) {
	Class := context.Query("class")
	questionList := con.AnswerService.GetQuestionsByClass(Class)
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   questionList,
			"code":   http.StatusOK,
		},
	)
}

// 根据科目和问题号获取答案
func (con *Controller) GetAnswerByClassAndQuestionId(context *gin.Context) {
	Class := context.Query("class")
	QuestionId, _ := strconv.Atoi(context.Query("questionId"))
	answerList := con.AnswerService.GetAnswerByClassAndQuestionId(Class, QuestionId)
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   answerList,
			"code":   http.StatusOK,
		},
	)
}

// 发送自测答案
func (con *Controller) SendTestAnswer(context *gin.Context) {
	testAnswers := make([]model.Answer, 0)
	if err := context.Bind(testAnswers); err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": "失败",
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
	}

}

// 获取科目列表
func (con *Controller) GetClassList(context *gin.Context) {
	classList := con.AnswerService.GetClassList()
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   classList,
			"code":   http.StatusOK,
		},
	)
}

// 新增科目
func (con *Controller) AddClass(context *gin.Context) {
	newClass := util.Class{
		ClassCode: strings.ToUpper(context.Query("classCode")),
		ClassName: context.Query("className"),
	}
	con.AnswerService.AddClass(newClass)
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   "",
			"code":   http.StatusOK,
		},
	)
}

// 根据科目code删除科目
func (con *Controller) DeleteClassByCode(context *gin.Context) {
	con.AnswerService.DeleteClass(context.Query("classCode"))
	context.JSON(
		http.StatusOK,
		gin.H{
			"result": "成功",
			"data":   "",
			"code":   http.StatusOK,
		},
	)
}

func (con *Controller) AutoSyncAnswer(context *gin.Context) {
	var user util.Itt
	var err error
	_ = context.Bind(&user)
	user.UserToken, err = user.LoginIttByuserNameAndPassword()
	if err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": err.Error(),
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
		return
	}
	fmt.Println("Token: ", user.UserToken)
	var ClassList []util.ClassList
	ClassList, err = user.GetClassList()
	if err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": err.Error(),
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
		return
	}
	fmt.Println("Class: ", ClassList)
	var msg string
	for _, c := range ClassList {
		var answerList []util.Answer
		answerList, err = user.GetAnswer(c.SkillEvaluationId)
		if err != nil {
			context.JSON(
				http.StatusOK,
				gin.H{
					"result": err.Error(),
					"data":   "",
					"code":   http.StatusInternalServerError,
				},
			)
			return
		}
		fmt.Println("Answer: ", answerList)
		err = con.AnswerService.JsonParseByJson(answerList, c.StackName)
		if err != nil {
			msg = msg + err.Error() + "\n"
		}
	}
	if msg == "" {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": user.UserName + "答案同步成功",
				"data":   "",
				"code":   http.StatusOK,
			},
		)
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{
				"result": msg,
				"data":   "",
				"code":   http.StatusInternalServerError,
			},
		)
	}

}
