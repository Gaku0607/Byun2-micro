package models

import (
	"encoding/json"
	"errors"
)

// 1000...1999 為Member_srv的錯誤
const (
	ERROR_MEMBER_NOTEXISTS         int = 1000 + iota //用戶不存在
	ERROR_PWD_NOTEXISTS                              //密碼不存在
	ERROR_MEMBER_EXISTS                              //用戶已存在
	ERROR_PWD_EXISTS                                 //密碼已存在
	ERROR_EXPIRED_CERTIFICATION_ID                   // 驗證碼過期
	ERROR_EMAIL_EXISTS                               //帳號已存在
)

//5000...5999為通用錯誤
const (
	ERRIR_VERSION_FAILD int = 5000 + iota //版本錯誤
	ERROR_NOSUCH_DATA                     //查無資料
	ERROR_SERVER_FAILD                    //伺服器錯誤
)

var NilErr = errors.New("nil")

type SrvErr struct {
	Stack string `json:"status"`
	Code  int    `json:"code"`
}

func NewErr(code int, err error) *SrvErr {
	return &SrvErr{
		Code:  code,
		Stack: err.Error(),
	}
}

//返回給用戶的錯誤grpc將返回 json格式
func (e *SrvErr) Error() string {
	obj, _ := json.Marshal(e)
	return string(obj)
}
