package service

import (
	"JsonParse/model"
	"JsonParse/repository"
	"JsonParse/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	paths "path"
	"path/filepath"
	"strings"
)

type AnswerService struct {
	AnswerRepo *repository.AnswerRepo
}

func (s *AnswerService) InitAnswer() {
	//数据文件清单
	fileList := make([]string, 0)
	if err := filepath.Walk("resource/", func(path string, info fs.FileInfo, err error) error {
		fileExt := ".txt.json."
		if !info.IsDir() && paths.Ext(info.Name()) != "" && strings.Contains(fileExt, paths.Ext(info.Name())) {
			fileList = append(fileList, info.Name())
		}
		return nil
	}); err != nil {
		log.Panicln(err)
		return
	}
	for _, fileName := range fileList {
		jsonFile, err := os.Open(fmt.Sprint("resource/", fileName))
		if err != nil {
			fmt.Println("open json file has error,please check: ", err.Error())
			continue
		}
		defer jsonFile.Close()
		decoder := json.NewDecoder(jsonFile)
		var answers []util.Answer
		class := strings.ToUpper(strings.Split(fileName, ".")[0])
		err = decoder.Decode(&answers)
		if err != nil {
			fmt.Println("read file has error,please check : ", err.Error())
			continue
		}
		answerList := make([]model.Answer, 0)
		for _, a := range answers {
			optioins, result := a.Question.OptionsToString()
			answerList = append(answerList, model.Answer{
				Class:           class,
				QuestionId:      a.Question.Id,
				Stem:            a.Question.Stem,
				Analysis:        a.Question.Analysis,
				Options:         optioins,
				Result:          result,
				KnowledgePoints: a.Question.KnowledgeToString(),
			})
		}
		//fmt.Println(answerList)
		s.AnswerRepo.InitTable(class, answerList)
	}
}

func (s *AnswerService) GetAnswerByClass(class string) []model.Answer {
	return s.AnswerRepo.GetAnswerByClass(class)
}

func (s *AnswerService) JsonParse(file *multipart.File, fileName string) error {
	// 获取文件名，并创建新的文件存储
	decoder := json.NewDecoder(*file)
	var answers []util.Answer
	class := strings.ToUpper(strings.Split(fileName, ".")[0])
	if !s.AnswerRepo.ClassIsEnable(model.Class{ClassCode: class}) {
		return errors.New("当前无法上传科目名为" + class + "的答案，请联系管理员")
	}
	err := decoder.Decode(&answers)
	if err != nil {
		return err
	}
	answerList := make([]model.Answer, 0)
	for _, a := range answers {
		optioins, result := a.Question.OptionsToString()
		answerList = append(answerList, model.Answer{
			Class:           class,
			QuestionId:      a.Question.Id,
			Stem:            a.Question.Stem,
			Analysis:        a.Question.Analysis,
			Options:         optioins,
			Result:          result,
			KnowledgePoints: a.Question.KnowledgeToString(),
		})
	}
	s.AnswerRepo.AddAnswer(answerList)
	return nil
}

func (s *AnswerService) JsonParseByJson(answers []util.Answer, className string) error {
	// 获取文件名，并创建新的文件存储
	if !s.AnswerRepo.ClassIsEnable(model.Class{ClassName: className}) {
		return errors.New("当前无法上传科目名为" + className + "的答案，请联系管理员")
	}
	classCode := s.AnswerRepo.GetClassCodeByName(className)
	answerList := make([]model.Answer, 0)
	for _, a := range answers {
		optioins, result := a.Question.OptionsToString()
		answerList = append(answerList, model.Answer{
			Class:           classCode,
			QuestionId:      a.Question.Id,
			Stem:            a.Question.Stem,
			Analysis:        a.Question.Analysis,
			Options:         optioins,
			Result:          result,
			KnowledgePoints: a.Question.KnowledgeToString(),
		})
	}
	s.AnswerRepo.AddAnswer(answerList)
	return nil
}

func (s *AnswerService) GetQuestionsByClass(class string) []int {
	return s.AnswerRepo.GetQuestionsByClass(class)
}

func (s *AnswerService) GetAnswerByClassAndQuestionId(class string, questionId int) model.Answer {
	return s.AnswerRepo.GetAnswerByClassAndQuestionId(class, questionId)
}

func (s *AnswerService) GetClassList() []model.Class {
	return s.AnswerRepo.GetClass()
}

func (s *AnswerService) AddClass(class util.Class) {
	a := model.Class{ClassCode: class.ClassCode, ClassName: class.ClassName}
	s.AnswerRepo.AddClass(&a)
}

func (s *AnswerService) DeleteClass(classCode string) {
	a := model.Class{ClassCode: classCode}
	s.AnswerRepo.DeleteClass(&a)
}
