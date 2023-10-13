package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type Itt struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	UserToken string `json:"token"`
}

type LoginResult struct {
	Code        int       `json:"code"`
	Data        UserToken `json:"data"`
	Msg         string    `json:"msg"`
	ServiceCode string    `json:"service_code"`
	TraceId     string    `json:"trace_id"`
}

type ClassListResult struct {
	Code        int         `json:"code"`
	Data        []ClassList `json:"data"`
	Msg         string      `json:"msg"`
	ServiceCode string      `json:"service_code"`
	TraceId     string      `json:"trace_id"`
}

type AnswerResult struct {
	Code        int           `json:"code"`
	Data        AnswerDetails `json:"data"`
	Msg         string        `json:"msg"`
	ServiceCode string        `json:"service_code"`
	TraceId     string        `json:"trace_id"`
}

type AnswerDetails struct {
	Details []Answer `json:"details"`
}

type UserToken struct {
	Token string `json:"token"`
}

type ClassList struct {
	StackName         string `json:"stack_name"`
	SkillEvaluationId string `json:"skill_evaluation_id"`
}

func (i *Itt) LoginIttByuserNameAndPassword() (string, error) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["jobNumber"] = i.UserName
	data["password"] = i.Password
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", viper.GetString("itt.url.login"), bytes.NewReader(bytesData))
	req.Header.Add("companyId", viper.GetString("itt.companyId"))
	req.Header.Add("Accept", viper.GetString("itt.Accept"))
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(resp.Body)
	var res LoginResult
	_ = json.Unmarshal(body, &res)
	if res.Code != 200 {
		return "", errors.New(res.Msg)
	}
	return "Bearer " + res.Data.Token, nil
}

func (i *Itt) GetClassList() ([]ClassList, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", viper.GetString("itt.url.class"), nil)
	req.Header.Add("companyId", viper.GetString("itt.companyId"))
	req.Header.Add("Accept", viper.GetString("itt.Accept"))
	req.Header.Add("Authorization", i.UserToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	var res ClassListResult
	_ = json.Unmarshal(body, &res)
	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}
	return res.Data, nil
}

func (i *Itt) GetAnswer(SkillEvaluationId string) ([]Answer, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", viper.GetString("itt.url.answer")+"?skill_evaluation_id="+SkillEvaluationId, nil)
	req.Header.Add("companyId", viper.GetString("itt.companyId"))
	req.Header.Add("Accept", viper.GetString("itt.Accept"))
	req.Header.Add("Authorization", i.UserToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	var res AnswerResult
	_ = json.Unmarshal(body, &res)
	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}
	return res.Data.Details, nil
}
