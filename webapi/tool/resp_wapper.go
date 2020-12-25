package tool

import (
	"encoding/json"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/gin-gonic/gin"
)

type HandlerWapper func(*gin.Context) (*models.ResultData, error)

func ResultWapper(fn HandlerWapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := fn(c)

		if err != nil {
			if srverr, ok := err.(*models.SrvErr); ok {
				Failed(c, srverr)
				return
			}

			var e models.SrvErr

			if jerr := json.Unmarshal([]byte(err.Error()), &e); jerr != nil {
				Failed(c, models.NewErr(models.ERROR_SERVER_FAILD, err))
			} else {
				Failed(c, &e)
			}

		} else {
			Success(c, data)
		}
	}
}

//請求成功Resp
func Success(c *gin.Context, data *models.ResultData) {
	c.JSON(
		data.StatusCode,
		gin.H{"code": 0, "data": data.Data},
	)
}

//請求失敗Resp
func Failed(c *gin.Context, err *models.SrvErr) {
	c.Set("Error", err)
	c.JSON(
		err.StatusCode(),
		gin.H{"code": err.ErrorCode(), "data": err.Error()},
	)
}
