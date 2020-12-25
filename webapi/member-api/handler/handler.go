package handler

import (
	"fmt"
	"net/http"
	"strconv"

	member "github.com/Gaku0607/Byun2-micro/webapi/member-api/proto/pb"
	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func getsrv(c *gin.Context) member.MemberService {
	s, _ := c.Get("service")
	return s.(member.MemberService)
}

func getVersion() *member.Version {
	return &member.Version{ApiVersion: apiVersion}
}

func Registry(c *gin.Context) (*models.ResultData, error) {

	var m models.Registry

	if err := c.ShouldBindJSON(&m); err != nil {
		return nil, err
	}

	req := member.RegisteyReq{
		Version: getVersion(),
		MInfo:   &member.Member{},
	}

	req.MInfo.Email = m.Email
	req.MInfo.Name = m.Name
	req.MInfo.Password = m.Pwd

	srv := getsrv(c)

	if _, err := srv.Registey(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		return &models.ResultData{StatusCode: http.StatusOK}, nil
	}
}

func Login(c *gin.Context) (*models.ResultData, error) {

	var parms models.Login

	if err := c.ShouldBindJSON(&parms); err != nil {
		return nil, err
	}

	var req member.LoginReq

	req.Version = getVersion()
	req.MName = parms.Name
	req.Password = parms.Pwd

	srv := getsrv(c)

	if resp, err := srv.Login(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		c.Header("Authorization", resp.Token)
		return &models.ResultData{Data: resp.Token, StatusCode: http.StatusOK}, nil
	}
}

func CreatMember(c *gin.Context) (*models.ResultData, error) {

	var Code models.Code

	if err := c.ShouldBind(&Code); err != nil {
		return nil, err
	}

	var req member.CreateReq

	req.Code = Code.Code
	req.Version = getVersion()

	srv := getsrv(c)

	if _, err := srv.Create(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		return &models.ResultData{StatusCode: http.StatusCreated}, nil
	}
}

func DeleteMember(c *gin.Context) (*models.ResultData, error) {

	Id := c.GetInt64("MemberId")
	srv := getsrv(c)

	var req member.DeleteReq

	req.Version = getVersion()
	req.MId = Id

	if _, err := srv.Delete(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		return &models.ResultData{StatusCode: http.StatusOK}, nil
	}
}

func ModityMemberInfo(c *gin.Context) (*models.ResultData, error) {

	Id := c.GetInt64("MemberId")
	var m models.MotidyMember

	if err := c.ShouldBindJSON(&m); err != nil {
		return nil, err
	}

	srv := getsrv(c)

	var req member.ModityReq

	req = member.ModityReq{
		Version: getVersion(),
		MId:     Id,
		MInfo: &member.Member{
			Password: m.Pwd,
			Email:    m.Email,
		},
	}

	fmt.Println(req)

	if _, err := srv.Modity(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		return &models.ResultData{StatusCode: http.StatusOK}, nil
	}
}

func GetMemberInfo(c *gin.Context) (*models.ResultData, error) {

	name := c.Param("name")

	srv := getsrv(c)

	var req member.MemberInfoReq

	req.Version = getVersion()
	req.MName = name

	if resp, err := srv.MemberInfo(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		m := &models.Member{}
		m.Name = resp.MName
		m.Banlancer = resp.Balanc
		m.IsSeller = resp.IsSeller
		return &models.ResultData{StatusCode: http.StatusOK, Data: m}, nil
	}
}

func MemberList(c *gin.Context) (data *models.ResultData, err error) {

	srv := getsrv(c)

	var req member.MemberListReq

	req = member.MemberListReq{
		Version:  getVersion(),
		Category: &member.MemberListReq_Category{},
	}

	val := c.Request.URL.Query()

	isseller := val.Get("is_seller")

	if isseller != "" {
		req.Category.IsSeller, _ = strconv.ParseBool(isseller)
	}

	offset := val.Get("offset")

	if offset != "" {
		if Offset, _ := strconv.ParseInt(offset, 10, 32); Offset < 1 {
			req.Category.Offset = 1
		} else {
			req.Category.Offset = int32(Offset)
		}
	}

	if resp, err := srv.MemberList(c.Request.Context(), &req); err != nil {
		return nil, err
	} else {
		list := make([]*models.Member, len(resp.MList))
		for index := range list {
			list[index] = &models.Member{
				Name:      resp.MList[index].MName,
				IsSeller:  resp.MList[index].IsSeller,
				Banlancer: resp.MList[index].Balanc,
			}
		}
		return &models.ResultData{StatusCode: http.StatusOK, Data: list}, nil
	}
}
