package permission

import (
	"context"
	"fmt"
	"log"

	"github.com/mrhumster/web-server-gin/gen/go/permission"
	"github.com/mrhumster/web-server-gin/internal/service"
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
		Error:   fmt.Sprintf("Check permission err %s", err),
	}, nil
}

func (s *PermissionGRPCServer) AddPolicy(ctx context.Context, req *permission.AddPolicyRequest) (*permission.AddPolicyResponse, error) {
	log.Printf("⚠️ PermissionGRPCServer: AddPolicy: policy=%s, resource=%s, permission=%s",
		req.GetPolicy(), req.GetResorce(), req.GetPermission())
	added, err := s.permissionServer.AddPolicy(req.GetPolicy(), req.GetResorce(), req.GetPermission())
	log.Printf("⚠️ added=%t", added)
	return &permission.AddPolicyResponse{
		Added: added,
		Error: fmt.Sprintf("Add policy err %s", err),
	}, nil
}

func (s *PermissionGRPCServer) RemovePolicy(ctx context.Context, req *permission.RemovePolicyRequest) (*permission.RemovePolicyResponse, error) {
	removed, err := s.permissionServer.RemovePolicy(req.GetPolicy(), req.GetResource(), req.GetPermission())
	return &permission.RemovePolicyResponse{
		Removed: removed,
		Error:   fmt.Sprintf("Remove policy err %s", err),
	}, nil
}

func (s *PermissionGRPCServer) AddPolicyIfNotExists(ctx context.Context, req *permission.AddPolicyIfNotExistsRequest) (*permission.AddPolicyIfNotExistsResponse, error) {
	exists, err := s.permissionServer.AddPolicyIfNotExists(req.GetPolicy(), req.GetResource(), req.GetPermission())
	return &permission.AddPolicyIfNotExistsResponse{
		Exists: exists,
		Error:  fmt.Sprintf("Add policy if exists %s", err),
	}, nil
}
