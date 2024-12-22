package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	pb "code.emcdtech.com/b2b/endpoint/protocol/b2bEndpointAdmin"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Admin struct {
	clientSrv service.Client
	secretSrv service.Secret
	ipSrv     service.IP
	pb.UnimplementedEndpointAdminServiceServer
}

func NewAdmin(clientSrv service.Client, secretSrv service.Secret, ipSrv service.IP) *Admin {
	return &Admin{
		clientSrv: clientSrv,
		secretSrv: secretSrv,
		ipSrv:     ipSrv,
	}
}

func (a *Admin) AddClient(ctx context.Context, req *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	id, err := a.clientSrv.Add(ctx, req.Name)
	if err != nil {
		log.Error(ctx, "admin: addClient: add: %s", err.Error())
		return nil, fmt.Errorf("add client: %w", err)
	}
	return &pb.AddClientResponse{
		Id: id.String(),
	}, nil
}

func (a *Admin) GenerateKey(ctx context.Context, req *pb.GenerateKeyRequest) (*pb.GenerateKeyResponse, error) {
	clientID, err := uuid.Parse(req.ClientId)
	if err != nil {
		log.Error(ctx, "admin: generateKey: parse: %s", err.Error())
		return nil, fmt.Errorf("parse: %w", err)
	}

	keys, err := a.secretSrv.GenerateKeys(ctx, clientID)
	if err != nil {
		log.Error(ctx, "admin: generateKey: generateKeys: %s", err.Error())
		return nil, fmt.Errorf("generateKeys: %w", err)
	}

	return &pb.GenerateKeyResponse{
		ApiKey:    keys.ApiKey.String(),
		ApiSecret: keys.ApiSecret.String(),
	}, nil
}

func (a *Admin) GetActiveKeys(ctx context.Context, req *pb.GetActiveKeysRequest) (*pb.GetActiveKeysResponse, error) {
	clientID, err := uuid.Parse(req.ClientId)
	if err != nil {
		log.Error(ctx, "admin: getActiveKeys: parse clientID: %s", err.Error())
		return nil, fmt.Errorf("parse clientID: %w", err)
	}

	secrets, err := a.secretSrv.GetActiveKeys(ctx, clientID)
	if err != nil {
		log.Error(ctx, "admin: getActiveKeys: getActiveKeys: %s", err.Error())
		return nil, fmt.Errorf("getActiveKeys: %w", err)
	}

	pbSecrets := make([]*pb.Secret, len(secrets))
	for i, secret := range secrets {
		pbSecrets[i] = &pb.Secret{
			ApiKey:    secret.ApiKey.String(),
			CreatedAt: timestamppb.New(secret.CreatedAt),
			LastUsed:  timestamppb.New(secret.LastUsed),
		}
	}
	return &pb.GetActiveKeysResponse{
		Secrets: pbSecrets,
	}, nil
}

func (a *Admin) DeactivateKey(ctx context.Context, req *pb.DeactivateKeyRequest) (*pb.DeactivateKeyResponse, error) {
	clientID, err := uuid.Parse(req.ClientId)
	if err != nil {
		log.Error(ctx, "admin: deactivateKey: parse clientID: %s", err.Error())
		return nil, fmt.Errorf("parse clientID: %w", err)
	}
	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		log.Error(ctx, "admin: deactivateKey: parse apiKey: %s", err.Error())
		return nil, fmt.Errorf("parse apiKey: %w", err)
	}

	err = a.secretSrv.DeactivateKey(ctx, clientID, apiKey)
	if err != nil {
		log.Error(ctx, "admin: deactivateKey: deactivateKey: %s", err.Error())
		return nil, fmt.Errorf("deactivateKey: %w", err)
	}

	return &pb.DeactivateKeyResponse{}, nil
}

func (a *Admin) DeactivateAllKeys(ctx context.Context, req *pb.DeactivateAllKeysRequest) (*pb.DeactivateAllKeysResponse, error) {
	clientID, err := uuid.Parse(req.ClientId)
	if err != nil {
		log.Error(ctx, "admin: deactivateAllKeys: parse clientID: %s", err.Error())
		return nil, fmt.Errorf("parse clientID: %w", err)
	}

	err = a.secretSrv.DeactivateAllKeys(ctx, clientID)
	if err != nil {
		log.Error(ctx, "admin: deactivateAllKeys: deactivateAllKeys: %s", err.Error())
		return nil, fmt.Errorf("deactivateAllKeys: %w", err)
	}

	return &pb.DeactivateAllKeysResponse{}, nil
}

func (a *Admin) AddIPs(ctx context.Context, req *pb.AddIPsRequest) (*pb.AddIPsResponse, error) {
	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		log.Error(ctx, "admin: addIPs: parse apiKey: %s", err.Error())
		return nil, fmt.Errorf("addIPs: parse apiKey: %w", err)
	}

	var (
		ips      []*model.IP
		parsedIP net.IP
	)
	for _, ip := range req.Ips {
		parsedIP = net.ParseIP(ip)
		if parsedIP == nil {
			log.Error(ctx, "admin: addIPs: parse ip: %s", ip)
			return nil, fmt.Errorf("addIPs: parse ip: %w", err)
		}
		ips = append(ips, &model.IP{
			ApiKey:    apiKey,
			Address:   ip,
			CreatedAt: time.Now().UTC(),
		})
	}
	if len(ips) == 0 {
		return &pb.AddIPsResponse{}, fmt.Errorf("addIPs: no ips found")
	}

	err = a.ipSrv.AddIPs(ctx, ips)
	if err != nil {
		log.Error(ctx, "admin: addIP: addIPs: %s", err.Error())
		return nil, fmt.Errorf("addIPs: %w", err)
	}

	return &pb.AddIPsResponse{}, nil
}

func (a *Admin) GetIPs(ctx context.Context, req *pb.GetIPsRequest) (*pb.GetIPsResponse, error) {
	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		log.Error(ctx, "admin: getIP: parse apiKey: %s", err.Error())
		return nil, fmt.Errorf("getIP: parse apiKey: %w", err)
	}

	ips, err := a.ipSrv.GetIPs(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("getIPs: %w", err)
	}

	pbIPs := make([]string, len(ips))

	for i, ip := range ips {
		pbIPs[i] = ip.Address
	}

	return &pb.GetIPsResponse{
		Ips: pbIPs,
	}, nil
}

func (a *Admin) DeleteIP(ctx context.Context, req *pb.DeleteIPRequest) (*pb.DeleteIPResponse, error) {
	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		log.Error(ctx, "admin: deleteIP: parse apiKey: %s", err.Error())
		return nil, fmt.Errorf("deleteIP: parse apiKey: %w", err)
	}

	ip := net.ParseIP(req.Ip)
	if ip == nil {
		log.Error(ctx, "admin: deleteIP: parse ip: %s", req.Ip)
		return nil, fmt.Errorf("deleteIP: parse ip: %w", err)
	}

	err = a.ipSrv.DeleteIP(ctx, apiKey, req.Ip)
	if err != nil {
		log.Error(ctx, "admin: deleteIP: deleteIP: %s", err.Error())
		return nil, fmt.Errorf("deleteIP: %w", err)
	}

	return &pb.DeleteIPResponse{}, nil
}

func (a *Admin) DeleteAllIPs(ctx context.Context, req *pb.DeleteAllIPsRequest) (*pb.DeleteAllIPsResponse, error) {
	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		log.Error(ctx, "admin: deleteAllIP: parse apiKey: %s", err.Error())
		return nil, fmt.Errorf("deleteAllIPs: %w", err)
	}

	err = a.ipSrv.DeleteAllIPs(ctx, apiKey)
	if err != nil {
		log.Error(ctx, "admin: deleteAllIP: deleteAllIPs: %s", err.Error())
		return nil, fmt.Errorf("deleteAllIPs: %w", err)
	}

	return &pb.DeleteAllIPsResponse{}, nil
}
