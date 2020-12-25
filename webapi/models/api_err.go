package models

import (
	"errors"
	"net/http"
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

//4000...4999為jwt錯誤
const (
	ERROR_INSUFFICIENT_PERMISSIONS int = 4001 + iota //權限不足
	ERROR_FORMAT_JWT                                 //Jwt 格式錯誤
	ERROR_EXPIRED_JWT                                //Jwt 過期
)

//5000...5999為通用錯誤
const (
	ERRIR_VERSION_FAILD int = 5000 + iota //版本錯誤
	ERROR_NOSUCH_DATA                     //查無資料
	ERROR_SERVER_FAILD                    //伺服器錯誤
)

//返回前端錯誤
var errmsg map[int]string = map[int]string{
	ERROR_MEMBER_NOTEXISTS:         "User not exist",
	ERROR_PWD_NOTEXISTS:            "Password  not exist",
	ERROR_MEMBER_EXISTS:            "User already exist",
	ERROR_PWD_EXISTS:               "Password already exist",
	ERROR_EMAIL_EXISTS:             "Email already exist",
	ERROR_EXPIRED_CERTIFICATION_ID: "Expired CertificationID Time",
	ERRIR_VERSION_FAILD:            "version failed",
	ERROR_INSUFFICIENT_PERMISSIONS: "Insufficient Permissions",
	ERROR_EXPIRED_JWT:              "Expired JWT Time",
	ERROR_FORMAT_JWT:               "JWT format failed",
	ERROR_SERVER_FAILD:             "server failed",
}

var statuscode map[int]int = map[int]int{
	ERROR_MEMBER_NOTEXISTS:         http.StatusBadRequest,
	ERROR_PWD_NOTEXISTS:            http.StatusUnauthorized,
	ERROR_MEMBER_EXISTS:            http.StatusBadRequest,
	ERROR_PWD_EXISTS:               http.StatusBadRequest,
	ERROR_EMAIL_EXISTS:             http.StatusBadRequest,
	ERROR_EXPIRED_CERTIFICATION_ID: http.StatusNotFound,
	ERRIR_VERSION_FAILD:            http.StatusNotFound,
	ERROR_INSUFFICIENT_PERMISSIONS: 0,
	ERROR_FORMAT_JWT:               0,
	ERROR_EXPIRED_JWT:              0,
	ERROR_SERVER_FAILD:             http.StatusInternalServerError,
}

var NilErr = errors.New("nil")

type SrvErr struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func NewErr(code int, err error) *SrvErr {
	return &SrvErr{
		Code:   code,
		Status: err.Error(),
	}
}

//返回給用戶的錯誤
func (e *SrvErr) Error() string {
	return errmsg[e.Code]
}

//自定義協議碼
func (e *SrvErr) ErrorCode() int {
	return e.Code
}

//http狀態碼
func (e *SrvErr) StatusCode() int {
	return statuscode[e.Code]
}

//系統錯誤
func (e *SrvErr) Stack() string {
	return e.Status
}
