package membermidd

import (
	"context"
	"os"

	member "github.com/Gaku0607/Byun2-micro/webapi/member-api/proto/pb"
	"github.com/Gaku0607/Byun2-micro/webapi/tool/banlancer"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

func MemberMidd() gin.HandlerFunc {

	srvname := os.Getenv("member_srv_name")

	var servicenode *registry.Node

	fn := func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			servicenode = node
			return cf(ctx, node, req, rsp, opts)
		}
	}

	b := memberBanlancer()

	service := micro.NewService(
		micro.Name(os.Getenv("api_name")+".client"),
		micro.Selector(b),
		micro.WrapCall(fn),
		micro.WrapClient(hystrix()),
	)

	membersrv := member.NewMemberService(
		srvname,
		service.Client(),
	)

	return func(c *gin.Context) {
		b.Init(
			selector.SetStrategy(banlancer.IPHash(c.ClientIP())),
		)

		c.Set("service", membersrv)

		c.Next()

		c.Set("servicenode", servicenode)
	}
}

func hystrix() client.Wrapper {
	return func(c client.Client) client.Client {
		return c
	}
}
