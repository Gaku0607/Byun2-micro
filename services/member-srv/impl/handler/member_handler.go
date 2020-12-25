package handler

import (
	"context"

	"github.com/Gaku0607/Byun2-micro/services/member-srv/impl/service"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/models"
	"github.com/Gaku0607/Byun2-micro/services/member-srv/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const apiVersion = "v1"

type MemberHandler struct {
	srv service.IMemberSrv
}

func NewMemberHandler(s service.IMemberSrv) *MemberHandler {
	return &MemberHandler{s}
}

func checkAPI(version string) error {
	if version != apiVersion {
		return models.NewErr(models.ERRIR_VERSION_FAILD, models.NilErr)
	}
	return nil
}

func (m *MemberHandler) Registey(ctx context.Context, req *pb.RegisteyReq, e *emptypb.Empty) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}
	return m.srv.RegisterVerification(req.MInfo)
}

func (m *MemberHandler) Create(ctx context.Context, req *pb.CreateReq, e *emptypb.Empty) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}
	return m.srv.CreateMember(req.Code)
}

func (m *MemberHandler) Delete(ctx context.Context, req *pb.DeleteReq, e *emptypb.Empty) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}
	return m.srv.DeleteMember(req.MId)
}

func (m *MemberHandler) Modity(ctx context.Context, req *pb.ModityReq, e *emptypb.Empty) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}
	return m.srv.MotidyMember(req.MId, req.MInfo)
}

func (m *MemberHandler) Login(ctx context.Context, req *pb.LoginReq, resp *pb.LoginResp) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}

	if resp.Token, err = m.srv.CheckLoginMsg(req.MName, req.Password); err != nil {
		return
	}
	return
}

func (m *MemberHandler) MemberList(ctx context.Context, req *pb.MemberListReq, resp *pb.MemberListResp) (err error) {

	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}
	list, err := m.srv.MemberList(req.Category)
	if err != nil {
		return
	}

	resp.MList = make([]*pb.MemberListResp_MemberList, len(list))

	for index, mb := range list {
		resp.MList[index] = &pb.MemberListResp_MemberList{
			MName:    mb.Name,
			IsSeller: *mb.IsSeller,
			Balanc:   mb.Balanc,
		}
	}

	return nil
}

func (m *MemberHandler) MemberInfo(ctx context.Context, req *pb.MemberInfoReq, resp *pb.MemberInfoResp) (err error) {
	if err = checkAPI(req.Version.ApiVersion); err != nil {
		return
	}

	member, err := m.srv.GetMemberInfo(req.MName)
	if err != nil {
		return
	}

	resp.MName = member.Name
	resp.IsSeller = *member.IsSeller
	resp.Balanc = member.Balanc

	return nil
}
