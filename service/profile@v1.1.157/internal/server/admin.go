package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
)

type Admin struct {
	profileService service.Profile
	log            service.ProfileLog
	pb.UnimplementedAdminProfileServiceServer
}

func NewAdmin(profileSrv service.Profile, profileLogSrv service.ProfileLog) *Admin {
	return &Admin{
		profileService: profileSrv,
		log:            profileLogSrv,
	}
}

const (
	deleteSub = "delete sub user"
	other     = "other action"
)

func (p *Admin) DeleteSubUser(ctx context.Context, req *pb.DeleteSubUserRequest) (*emptypb.Empty, error) {
	const op = "server.Admin.DeleteSubUser"
	resp := &emptypb.Empty{}
	newParentID, err := uuid.Parse(req.NewParentId)
	if err != nil {
		log.Error(ctx, "%s: parsing parent id: %s: %v", op, req.NewParentId, err)
		return resp, fmt.Errorf("%s: parse parent id", op)
	}

	userID, err := uuid.Parse(req.SubuserId)
	if err != nil {
		log.Error(ctx, "%s: parsing subuser id: %s: %v", op, req.SubuserId, err)
		return resp, fmt.Errorf("%s: parse subuser id", op)
	}

	err = p.profileService.SoftDeleteSubUser(context.WithoutCancel(ctx), userID, newParentID)
	if err != nil {
		log.Error(ctx, "%s: parent: %s subuser: %s: %v", op, newParentID, userID, err)
		return resp, fmt.Errorf("%s: soft delete", op)
	}

	dataMap := map[string]any{
		"action":      "soft subaccount delete",
		"subuser_id":  userID,
		"new parent":  newParentID,
		"modified by": "admin",
	}

	data, err := toString(dataMap)
	if err != nil {
		log.Error(ctx, "%s: fail to prepare log : %v", op, err)
		return &emptypb.Empty{}, nil
	}

	info := &model.ProfileLog{
		ChangeType: deleteSub,
		Details:    data,
		Originator: "admin",
	}

	if err := p.log.Log(ctx, info); err != nil {
		log.Error(ctx, "%s: fail to log changes: %v", op, err)
	}

	return &emptypb.Empty{}, nil
}

func toString(data map[string]any) (string, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
