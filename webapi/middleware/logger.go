package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	type WebLogger struct {
		UserAgent     string      `json:"user_agent"`             //請求
		SpendTime     int64       `json:"spend_time"`             //請求總花費時間
		ClientIP      string      `json:"client_ip"`              //客戶端ＩＰ
		HostName      string      `json:"host_name"`              //主機名稱
		StatusCode    int         `json:"status_code"`            //http狀態碼
		RequestMethod string      `json:"req_method"`             //請求方法
		RequestURL    string      `json:"req_url"`                //請求地址
		DataSize      int         `json:"data_size"`              //資料大小
		ServiceNode   interface{} `json:"service_node,omitempty"` //請求所調用的Node
		ErrCode       int         `json:"err_code,omitempty"`     //自定義錯誤碼
		ErrMsg        string      `json:"err_msg,omitempty"`      //錯誤信息
	}
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger := &WebLogger{}
		logger.SpendTime = time.Now().Sub(start).Milliseconds()
		logger.ClientIP = c.ClientIP()
		logger.RequestMethod = c.Request.Method
		logger.RequestURL = c.Request.URL.String()
		logger.DataSize = c.Writer.Size()
		logger.UserAgent = c.Request.UserAgent()
		logger.StatusCode = c.Writer.Status()

		name, err := os.Hostname()

		if err != nil {
			name = "unknow"
		}

		logger.HostName = name

		node, _ := c.Get("servicenode")

		logger.ServiceNode = node

		e, exits := c.Get("Error")

		if exits {
			if err, ok := e.(*models.SrvErr); ok {
				logger.ErrMsg = err.Stack()
				logger.ErrCode = err.ErrorCode()
			} else {
				logger.ErrMsg = e.(error).Error()
			}
		}

		data, _ := json.Marshal(&logger)

		fmt.Println(string(data))
	}
}
