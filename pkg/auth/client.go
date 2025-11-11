package auth

import "github.com/mrhumster/web-server-gin/internal/permission"

func NewPermissionClient(url string) (*permission.PermissionGRPCClient, error) {
	return permission.NewPermissionGRPCClient(url)
}
