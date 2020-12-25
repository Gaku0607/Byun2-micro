package service

import (
	"log"
	"testing"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/repository"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
)

var i IMemberSrv

func TestMain(m *testing.M) {
	i = NewMemberDataSrv(repository.NewMemberStoreInSQL())
	if err := i.Init(); err != nil {
		log.Println(err.Error())
		return
	}
	m.Run()
}

func TestLogin(t *testing.T) {
	token, err := i.CheckLoginMsg("dader", "1234")

	if err != nil {
		t.Error(err.Error())
	}

	println(token)
}

func TestCreateMember(t *testing.T) {
	//輸入驗證碼 驗證碼必須從TestRegister 獲取
	if err := i.CreateMember("4olizih1"); err != nil {
		t.Error(err.Error())
	}
}

func TestRegister(t *testing.T) {
	m := &pb.Member{}
	m.Name = "dader"
	m.Password = "1234"
	m.Email = "testemail"

	if err := i.RegisterVerification(m); err != nil {
		t.Error(err.Error())
	}
}
