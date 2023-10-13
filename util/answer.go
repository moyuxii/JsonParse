package util

import "fmt"

type Answer struct {
	Question Question `json:"question"`
}

type Question struct {
	Stem             string    `json:"stem"`
	Id               int       `json:"id"`
	Options          []Options `json:"options"`
	Analysis         string    `json:"analysis"`
	Knowledge_points []int     `json:"knowledge_point_ids"`
}

type Options struct {
	Correct bool   `json:"correct"`
	Id      int    `json:"id"`
	Option  string `json:"option"`
}

type Class struct {
	ClassCode string `json:"classCode"`
	ClassName string `json:"className"`
}

func (q *Question) OptionsToString() (string, int) {
	var questions string
	var result int
	for j, option := range q.Options {
		if option.Correct {
			result = j + 1
		}
	}
	return questions, result
}

func (q *Question) KnowledgeToString() string {
	var knowledge string
	for _, k := range q.Knowledge_points {
		knowledge = fmt.Sprint(knowledge, ",", k)
	}
	return knowledge
}
