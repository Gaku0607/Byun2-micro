package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/handler"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/repository"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/service"
	member "github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
	"github.com/joho/godotenv"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

var (
	registryaddress string
	servicename     string
	ttl             time.Duration
	interval        time.Duration
)

func Init() (err error) {
	if err = godotenv.Load("member_srv.env"); err != nil {
		return err
	}

	if registryaddress = os.Getenv("registry_address"); registryaddress == "" {
		return errors.New("registry_address is empty")
	}

	if servicename = os.Getenv("service_name"); servicename == "" {
		return errors.New("service_name is empty")
	}

	if ttlstr := os.Getenv("ttl"); ttlstr == "" {
		return errors.New("weights is empty")
	} else {
		if ttlint, err := strconv.Atoi(ttlstr); err != nil {
			return err
		} else {
			ttl = time.Duration(ttlint) * time.Second
		}
	}

	if intervalstr := os.Getenv("interval"); intervalstr == "" {
		return errors.New("interval is empty")
	} else {
		if intervalint, err := strconv.Atoi(intervalstr); err != nil {
			return err
		} else {
			interval = time.Duration(intervalint) * time.Second
		}
	}

	return

}
func main() {
	if err := Init(); err != nil {
		log.Fatalln(err.Error())
	}

	membersrv := service.NewMemberDataSrv(repository.NewMemberStoreInSQL())

	defer membersrv.Close()

	if err := membersrv.Init(); err != nil {
		log.Fatalln(err.Error())
	}

	r := consul.NewRegistry(
		registry.Addrs(registryaddress),
	)

	srv := micro.NewService(
		micro.Registry(r),
		micro.Name(servicename),
		micro.RegisterTTL(ttl),
		micro.RegisterInterval(interval),
	)

	member.RegisterMemberServiceHandler(
		srv.Server(),
		handler.NewMemberHandler(membersrv),
	)

	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
