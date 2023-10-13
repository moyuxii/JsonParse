package model

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	Class           string `json:"class" gorm:"class"`
	QuestionId      int    `json:"questionId" gorm:"questioin_id"`
	Stem            string `json:"stem" gorm:"stem"`
	Analysis        string `json:"analysis" gorm:"analysis"`
	Options         string `json:"options" gorm:"options"`
	Result          int    `json:"result" gorm:"result"`
	KnowledgePoints string `json:"knowledgePoints" gorm:"knowledge_points"`
}

type Class struct {
	gorm.Model
	ClassName string `json:"className" gorm:"class_name"`
	ClassCode string `json:"classCode" gorm:"class_code"`
}
