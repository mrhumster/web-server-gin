package permission

import (
	"context"
	"fmt"
	"log"

	"github.com/mrhumster/web-server-gin/proto/gen/go/permission"
	"github.com/mrhumster/web-server-gin/service"
)

type PermissionGRPCServer struct {
	permission.UnimplementedPermissionServiceServer
	permissionServer *service.PermissionService
}

func NewPermissionGRPCServer(p *service.PermissionService) *PermissionGRPCServer {
	if p == nil {
		log.Fatal("❌ PermissionService is nil in NewPermissionGRPCServer")
	}
	return &PermissionGRPCServer{
		permissionServer: p,
	}
}

func (s *PermissionGRPCServer) CheckPermission(ctx context.Context, req *permission.CheckPermissionRequest) (*permission.CheckPermissionResponse, error) {
	log.Printf("⚠️ PermissionGRPCServer: CheckPermission: user_id=%s, resource=%s, action=%s",
		req.GetUserId(), req.GetResource(), req.GetAction())

	allowed, err := s.permissionServer.CheckPermission(req.GetUserId(), req.GetResource(), req.GetAction())
	log.Printf("⚠️ allowed=%t", allowed)
	return &permission.CheckPermissionResponse{
		Allowed: allowed,
		Error:   fmt.Sprintf("Permission err %s", err),
	}, nil
}
