package middleware

import (
	"fmt"
	"time"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/Gaku0607/Byun2-micro/webapi/tool"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const deadline int64 = 10800

func MemberJwtMidd() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			tool.Failed(c, models.NewErr(models.ERROR_INSUFFICIENT_PERMISSIONS, models.NilErr))
			c.Abort()
			return
		}

		tokenclaims, err := jwt.ParseWithClaims(
			auth,
			&models.Claims{},
			func(t *jwt.Token) (interface{}, error) {
				return tool.Key, nil
			},
		)

		if err != nil {

			result := err.(*jwt.ValidationError)
			//如果jwt錯誤不為為過期
			if result.Errors != jwt.ValidationErrorExpired {
				tool.Failed(c, models.NewErr(models.ERROR_FORMAT_JWT, err))
				c.Abort()
				return
			}
			Claims := tokenclaims.Claims.(*models.Claims)
			//過期時間大餘3小時錯誤
			if time.Now().Unix() > Claims.ExpiresAt+deadline {
				tool.Failed(c, models.NewErr(models.ERROR_EXPIRED_JWT, err))
				c.Abort()
				return
			}
			//重發jwt
			token, err := tool.ResendJwt(Claims)
			if err != nil {
				tool.Failed(c, err.(*models.SrvErr))
				c.Abort()
				return
			}

			fmt.Println(Claims)

			c.Header("Authorization", token)
			c.Set("IsSeller", Claims.IsSeller)
			c.Set("MemberId", int64(Claims.MemeberId))
			c.Set("MemberName", Claims.MemberName)
			c.Next()
			return

		}
		Claims := tokenclaims.Claims.(*models.Claims)

		c.Set("IsSeller", Claims.IsSeller)
		c.Set("MemberId", int64(Claims.MemeberId))
		c.Set("MemberName", Claims.MemberName)
		c.Next()
	}
}
