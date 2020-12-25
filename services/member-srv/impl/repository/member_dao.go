package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/models"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
	"github.com/jinzhu/gorm"
)

var defualtsize int32 = 10

type MemberStoreInSQL struct {
	*gorm.DB
	pagesize int32
}

func NewMemberStoreInSQL() MemberRepository {
	return &MemberStoreInSQL{}
}

//繼承IMemberRepository的Init
func (r *MemberStoreInSQL) Init() error {
	conf, err := repositoryConfig()
	if err != nil {
		return err
	}

	if r.DB, err = repositoryConnect(conf); err != nil {
		return err
	}

	r.pagesize = defualtsize
	return nil
}

//繼承IMemberRepository的InitTable
func (r *MemberStoreInSQL) InitTable() error {
	if !r.HasTable(&models.Member{}) {
		if err := r.CreateTable(&models.Member{}).Error; err != nil {
			return err
		}
	}
	if !r.HasTable(&models.Email{}) {
		if err := r.CreateTable(&models.Email{}).Error; err != nil {
			return err
		}
	}
	return nil
}

//繼承IMemberRepository的MemberList
func (r *MemberStoreInSQL) MemberList(c *pb.MemberListReq_Category) ([]*models.Member, error) {
	list := make([]*models.Member, r.pagesize)

	db := r.DB.Model(&models.Member{}).Select("name,is_seller,balanc")

	if c.IsSeller {
		db = db.Where("is_seller = ?", 1)
	}

	return list, errResult(db.Scopes(r.PageLimit(c.Offset, r.pagesize)).Find(&list).Error)
}

//繼承IMemberRepository的QueryMemberByName
func (r *MemberStoreInSQL) QueryMemberByName(name string) (member *models.Member, err error) {
	member = &models.Member{}
	return member, r.Where("name = ?", name).First(member).Error
}

//繼承IMemberRepository的QueryMemberById
func (r *MemberStoreInSQL) QueryMemberById(Id int64) (member *models.Member, err error) {
	member = &models.Member{}
	return member, errResult(r.First(member, Id).Error)
}

//繼承IMemberRepository的CheckRegisterInfo
func (r *MemberStoreInSQL) CheckRegisterInfo(m *pb.Member) error {
	// if !r.checkRepectPwd(m.Password) {
	// 	return models.NewErr(models.ERROR_PWD_EXISTS, models.NilErr)
	// }
	if !r.checkRepectEmail(m.Email) {
		return models.NewErr(models.ERROR_EMAIL_EXISTS, models.NilErr)
	}
	if !r.checkRepectName(m.Name) {
		return models.NewErr(models.ERROR_MEMBER_EXISTS, models.NilErr)
	}
	return nil
}

//確認是否有重複密碼
// func (r *MemberStoreInSQL) checkRepectPwd(pwd string) bool {
// 	return r.Model(&models.Member{}).Where("pwd = ?", pwd).First(&models.Member{}).RecordNotFound()
// }

//確認是否有無重複姓名
func (r *MemberStoreInSQL) checkRepectName(name string) bool {
	return r.Model(&models.Member{}).Where("name = ?", name).First(&models.Member{}).RecordNotFound()
}

//確認是否有無重複信箱
func (r *MemberStoreInSQL) checkRepectEmail(email string) bool {
	return r.Model(&models.Email{}).Where("addrs = ?", email).First(&models.Email{}).RecordNotFound()
}

//建立用戶
func (r *MemberStoreInSQL) CreateMember(m *pb.Member) (err error) {
	member := &models.Member{
		Name: m.Name,
		Pwd:  m.Password,
		Email: models.Email{
			Addrs: m.Email,
		},
	}

	if err = r.Create(member).Error; err != nil {
		if strings.HasPrefix(fmt.Sprintf("%s", err), "Error 1062: Duplicate entry") {
			return r.CheckRegisterInfo(m)
		} else {
			return models.NewErr(models.ERROR_SERVER_FAILD, err)
		}
	}
	return nil
}

//繼承IMemberRepository的DeleteMemberById
func (r *MemberStoreInSQL) DeleteMemberById(Id int64) (err error) {
	tx := r.Begin()

	if err := r.delMember(tx, Id); err != nil {
		tx.Rollback()
		return err
	}
	if err := r.delEmail(tx, Id); err != nil {
		tx.Rollback()
		return err
	}
	return errResult(tx.Commit().Error)
}

//刪除會員
func (r *MemberStoreInSQL) delMember(tx *gorm.DB, Id int64) (err error) {
	if err = r.beforeDelMember(tx, Id); err != nil {
		return err
	}
	return errResult(tx.Delete(&models.Member{}, Id).Error)
}

//刪除用戶前必須調用（修改名字）
func (r *MemberStoreInSQL) beforeDelMember(tx *gorm.DB, Id int64) error {
	return errResult(tx.Model(&models.Member{}).Where("Id = ?", Id).UpdateColumn("name", gorm.Expr("concat(name,-?)", time.Now().Unix())).Error)
}

//刪除信箱
func (r *MemberStoreInSQL) delEmail(tx *gorm.DB, Id int64) (err error) {
	return errResult(tx.Where("member_id = ?", Id).Delete(&models.Email{}).Error)
}

//繼承IMemberRepository的ModityMember
func (r *MemberStoreInSQL) ModityMember(Id int64, m *pb.Member) (err error) {
	tx := r.Begin()

	if m.Email != "" {
		if err = r.modityEmail(tx, Id, m.Email); err != nil {
			tx.Rollback()
			fmt.Println(err)
			return err
		}
	}
	if m.Password != "" {
		if err = r.modityPwd(tx, Id, m.Password); err != nil {
			tx.Rollback()
			return err
		}
	}

	return errResult(tx.Commit().Error)
}

//修改密碼
func (r *MemberStoreInSQL) modityPwd(tx *gorm.DB, Id int64, pwd string) (err error) {
	return errResult(tx.Model(&models.Member{}).Where("id = ?", Id).Update("pwd", pwd).Error)
}

//修改信箱
func (r *MemberStoreInSQL) modityEmail(tx *gorm.DB, Id int64, addrs string) (err error) {
	return errResult(tx.Save(&models.Email{Addrs: addrs, MemberId: Id}).Error)
}

//繼承IMemberRepository的Close
func (r *MemberStoreInSQL) Close() error {
	return r.DB.Close()
}

//分頁
func (r *MemberStoreInSQL) PageLimit(offset, pagesize int32) func(*gorm.DB) *gorm.DB {
	return func(g *gorm.DB) *gorm.DB {
		return g.Offset((offset - 1) * pagesize).Limit(pagesize)
	}
}

//設置錯誤
func errResult(err error) error {
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.NewErr(models.ERROR_MEMBER_NOTEXISTS, err)
		} else {
			return models.NewErr(models.ERROR_SERVER_FAILD, err)
		}
	}
	return nil
}
