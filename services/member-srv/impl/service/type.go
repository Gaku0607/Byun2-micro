package service

import (
	"github.com/Gaku0607/Byun2-micro/services/member-srv/models"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
)

type IMemberSrv interface {
	//初始化參數
	Init() error
	//註冊驗證
	RegisterVerification(*pb.Member) error
	//確認登入信息
	CheckLoginMsg(string, string) (string, error)
	//建立用戶
	CreateMember(string) error
	//修改用戶
	MotidyMember(int64, *pb.Member) error
	//刪除用戶
	DeleteMember(int64) error
	//用戶列表
	MemberList(*pb.MemberListReq_Category) ([]*models.Member, error)
	//獲取用戶信息
	GetMemberInfo(string) (*models.Member, error)
	//關閉相關資源
	Close() error
}
