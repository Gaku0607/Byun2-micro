package membermidd

import (
	"os"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/Gaku0607/Byun2-micro/webapi/tool"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func configBanlancerParms() (address string, err error) {
	if address = os.Getenv("member_srv_registry_address"); address == "" {
		return "", models.ErrServiceRegAddressIsEmpty
	}
	return
}

func memberBanlancer() selector.Selector {
	address, err := configBanlancerParms()
	if err != nil {
		tool.ErrChan <- err
	}
	r := consul.NewRegistry(
		registry.Addrs(address),
	)

	s := selector.NewSelector(
		selector.Registry(r),
	)
	return s
}
