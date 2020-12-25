package service

import (
	"encoding/json"
	"fmt"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/repository"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/models"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/util"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

type MemberDataSrv struct {
	store repository.MemberRepository
	*redis.Pool
}

func NewMemberDataSrv(store repository.MemberRepository) IMemberSrv {
	return &MemberDataSrv{store: store}
}

//繼承IMemberSrv的Init
func (d *MemberDataSrv) Init() error {
	conf, err := configParms()
	if err != nil {
		return err
	}
	if err := d.store.Init(); err != nil {
		return err
	}
	if err := d.store.InitTable(); err != nil {
		return err
	}

	d.Pool, err = casheCommet(conf)
	return err
}

//繼承IMemberSrv的RegisterVerification
func (d *MemberDataSrv) RegisterVerification(m *pb.Member) error {

	if err := d.store.CheckRegisterInfo(m); err != nil {
		return err
	}

	code := util.GenerateCode()

	newCode, err := d.saveCashe(m, code)
	if err != nil {
		return err
	}

	return d.sendEmail(m.Email, newCode)
}

//儲存至緩存中
func (d *MemberDataSrv) saveCashe(m *pb.Member, code string) (newCode string, err error) {
	conn := d.Pool.Get()
	defer conn.Close()

	obj, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	for {
		result, err := conn.Do("setex", code, 600, obj)
		if err != nil {
			return "", models.NewErr(models.ERROR_SERVER_FAILD, err)
		}
		if result != 0 {
			return code, nil
		} else {
			code = util.GenerateCode()
		}
	}
}

//發送驗證碼至指定信箱
func (d *MemberDataSrv) sendEmail(addrs, code string) error {
	fmt.Println("addrs:", addrs, "code:", code)
	return nil
}

//刪除緩存中的code
func (d *MemberDataSrv) delCashe(code string) {
	conn := d.Pool.Get()
	defer conn.Close()
	conn.Do("del", code)
}

//繼承IMemberSrv的CreateMember
func (d *MemberDataSrv) CreateMember(code string) error {
	m, err := d.checkVerificationCode(code)
	if err != nil {
		return err
	}

	if m.Password, err = d.hashPwd(m.Password); err != nil {
		return err
	}

	return d.store.CreateMember(m)
}

//驗證驗證碼
func (d *MemberDataSrv) checkVerificationCode(code string) (m *pb.Member, err error) {
	conn := d.Pool.Get()
	defer conn.Close()

	memberobj, err := conn.Do("get", code)

	if err != nil {
		return nil, models.NewErr(models.ERROR_SERVER_FAILD, err)
	}

	if memberobj == nil {
		return nil, models.NewErr(models.ERROR_EXPIRED_CERTIFICATION_ID, models.NilErr)
	}

	d.delCashe(code)

	m = &pb.Member{}

	if err = json.Unmarshal(memberobj.([]byte), m); err != nil {
		return nil, models.NewErr(models.ERROR_SERVER_FAILD, err)
	}
	return m, nil
}

//繼承IMemberSrv的DeleteMember
func (d *MemberDataSrv) DeleteMember(Id int64) error {
	return d.store.DeleteMemberById(Id)
}

//繼承IMemberSrv的MotidyMember
func (d *MemberDataSrv) MotidyMember(Id int64, m *pb.Member) (err error) {
	if m.Password != "" {
		if m.Password, err = d.hashPwd(m.Password); err != nil {
			return err
		}
	}
	return d.store.ModityMember(Id, m)
}

//繼承IMemberSrv的MemberList
func (d *MemberDataSrv) MemberList(c *pb.MemberListReq_Category) ([]*models.Member, error) {
	return d.store.MemberList(c)
}

//繼承IMemberSrv的GetMemberInfo
func (d *MemberDataSrv) GetMemberInfo(name string) (member *models.Member, err error) {
	if member, err = d.store.QueryMemberByName(name); err != nil {
		return nil, models.NewErr(models.ERROR_MEMBER_NOTEXISTS, err)
	}
	return member, nil
}

//繼承IMemberSrv的CheckLoginMsg
func (d *MemberDataSrv) CheckLoginMsg(name string, pwd string) (token string, err error) {
	member, err := d.store.QueryMemberByName(name)
	if err != nil {
		return "", models.NewErr(models.ERROR_MEMBER_NOTEXISTS, err)
	}
	if err = d.isSamePwd(pwd, member.Pwd); err != nil {
		return "", err
	}
	//發放jwt
	return util.OutPutJwt(member)
}

//密碼是否相同
func (d *MemberDataSrv) isSamePwd(msgpwd, pwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(msgpwd)); err != nil {
		return models.NewErr(models.ERROR_PWD_NOTEXISTS, err)
	}
	return nil
}

//密碼哈希加密
func (d *MemberDataSrv) hashPwd(pwd string) (newpwd string, err error) {
	if bytepwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10); err != nil {
		return "", models.NewErr(models.ERROR_SERVER_FAILD, err)
	} else {
		return string(bytepwd), nil
	}
}

//關閉資源
func (d *MemberDataSrv) Close() error {
	d.store.Close()
	return d.Pool.Close()
}
