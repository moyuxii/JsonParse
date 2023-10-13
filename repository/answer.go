package repository

import (
	"JsonParse/model"
	"log"
)

type AnswerRepo struct {
	DB model.DataBase
}

func (answerRepo *AnswerRepo) InitTable(class string, answers []model.Answer) {
	if err := answerRepo.DB.SqlLite.Where("class = ?", class).Delete(&model.Answer{}).Error; err != nil {
		log.Panicln(err)
		return
	}
	answerRepo.DB.SqlLite.Save(answers)
}

func (answerRepo *AnswerRepo) AddAnswer(answers []model.Answer) {
	answerRepo.DB.SqlLite.Save(answers)
	answerRepo.DB.SqlLite.Where("id in (" +
		"select id from (select id,row_number() over(partition by class,question_id order by created_at desc) rm from answers where deleted_at is null) where rm <> 1" +
		")").Delete(&model.Answer{})
}

func (a *AnswerRepo) GetAnswerByClass(class string) []model.Answer {
	var answerList []model.Answer
	a.DB.SqlLite.Where("class = ?", class).Find(&answerList)
	return answerList
}

func (a *AnswerRepo) GetAnswerByClassAndQuestionId(class string, questionId int) model.Answer {
	var answer model.Answer
	a.DB.SqlLite.Where("class = ? and question_id = ? ", class, questionId).First(&answer)
	return answer
}

func (a *AnswerRepo) GetQuestionsByClass(class string) []int {
	var answerList []model.Answer
	a.DB.SqlLite.Distinct("question_id").Where("class = ?", class).Find(&answerList)
	var questionList []int
	for _, q := range answerList {
		questionList = append(questionList, q.QuestionId)
	}
	return questionList
}

func (a *AnswerRepo) GetClass() []model.Class {
	var classList []model.Class
	a.DB.SqlLite.Find(&classList)
	return classList
}

func (a *AnswerRepo) AddClass(class *model.Class) {
	a.DB.SqlLite.Save(&class)
}

func (a *AnswerRepo) DeleteClass(class *model.Class) {
	a.DB.SqlLite.Where(&class).Delete(&model.Class{})
	a.DB.SqlLite.Where("class = ? ", class.ClassCode).Delete(&model.Answer{})
}

func (a *AnswerRepo) ClassIsEnable(class model.Class) bool {
	var size int64
	a.DB.SqlLite.Model(model.Class{}).Where(&class).Count(&size)
	return size > 0
}

func (a *AnswerRepo) GetClassCodeByName(className string) string {
	var class model.Class
	a.DB.SqlLite.Where("class_name = ?", className).First(&class)
	return class.ClassCode
}
