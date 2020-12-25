package repository

import (
	"log"
	"testing"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
)

var s MemberRepository

func TestMain(m *testing.M) {
	s = NewMemberStoreInSQL()
	if err := s.Init(); err != nil {
		log.Println(err.Error())
		return
	}
	if err := s.InitTable(); err != nil {
		log.Println(err.Error())
		return
	}

	defer s.Close()

	m.Run()
}
func TestALL(t *testing.T) {
	t.Run("CreateMember", TestCreateMember)
	t.Run("MotidyMember", TestMotidyMember)
	t.Run("DelMember", TestDelMember)
}
func TestCreateMember(t *testing.T) {
	data := &pb.Member{
		Name:     "gaku",
		Password: "123",
		Email:    "124",
	}
	if err := s.CreateMember(data); err != nil {
		t.Error(err.Error())
	}
}

func TestDelMember(t *testing.T) {
	if err := s.DeleteMemberById(5); err != nil {
		t.Error(err.Error())
	}
}

func TestMotidyMember(t *testing.T) {
	data := &pb.Member{
		Password: "1233",
		Email:    "stt",
	}
	if err := s.ModityMember(1, data); err != nil {
		t.Error(err.Error())
	}
}
