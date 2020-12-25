package repository

import (
	"github.com/Gaku0607/Byun2-micro/services/member-srv/models"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
)

type MemberRepository interface {
	//初始化參數
	Init() error
	//初始化數據表
	InitTable() error
	//關閉數據庫資源
	Close() error
	//姓名查詢用戶
	QueryMemberByName(string) (*models.Member, error)
	//ID查詢用戶
	QueryMemberById(int64) (*models.Member, error)
	//獲取會員清單
	MemberList(*pb.MemberListReq_Category) ([]*models.Member, error)
	//確認註冊內容
	CheckRegisterInfo(*pb.Member) error
	//建立用戶
	CreateMember(m *pb.Member) error
	//ＩＤ刪除用戶
	DeleteMemberById(int64) error
	//修改用戶
	ModityMember(int64, *pb.Member) error
}
