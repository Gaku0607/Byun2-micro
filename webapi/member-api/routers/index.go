package routers

import (
	"net/http"

	"github.com/Gaku0607/Byun2-micro/webapi/member-api/handler"
	membermidd "github.com/Gaku0607/Byun2-micro/webapi/member-api/member-midd"
	"github.com/Gaku0607/Byun2-micro/webapi/middleware"
	"github.com/Gaku0607/Byun2-micro/webapi/tool"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	v1Group := engine.Group("/v1", middleware.Logger(), membermidd.MemberMidd())
	{
		//驗證註冊內容
		v1Group.POST("/registry", tool.ResultWapper(handler.Registry))
		//登入
		v1Group.POST("/login", tool.ResultWapper(handler.Login))
		//獲取用戶清單
		v1Group.GET("/members", tool.ResultWapper(handler.MemberList))
		//獲取用戶信息
		v1Group.GET("/members/:name", tool.ResultWapper(handler.GetMemberInfo))
		//驗證創建用戶
		v1Group.POST("/member", tool.ResultWapper(handler.CreatMember))

		//登入後
		MGroup := v1Group.Group("/", middleware.MemberJwtMidd())
		{
			//修改
			MGroup.PUT("/member", tool.ResultWapper(handler.ModityMemberInfo))
			//刪除
			MGroup.DELETE("/member", tool.ResultWapper(handler.DeleteMember))
		}
	}

	engine.NoRoute(
		func(c *gin.Context) {
			c.Status(http.StatusNotFound)
			c.String(http.StatusNotFound, "<h1>Error</h1>")
		},
	)

	return engine
}
