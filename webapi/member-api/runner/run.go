package runner

import (
	"errors"
	"time"

	"github.com/Gaku0607/Byun2-micro/webapi/member-api/routers"
	"github.com/Gaku0607/Byun2-micro/webapi/tool"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func Run() {

	var err error

	defer func() {
		if err == nil {
			tool.ErrChan <- errors.New("shutdown")
		}
	}()

	r := consul.NewRegistry(
		registry.Addrs(regaddrs),
	)

	service := web.NewService(
		web.Registry(r),
		web.Name(name),
		web.RegisterInterval(time.Duration(10)*time.Second),
		web.RegisterTTL(time.Duration(60)*time.Second),
		web.Address(port),
		web.Handler(routers.Routers()),
	)

	if err = service.Init(); err != nil {
		tool.ErrChan <- err
	}

	if err = service.Run(); err != nil {
		tool.ErrChan <- err
	}

}
