package auth

import "github.com/mrhumster/web-server-gin/internal/permission"

type PermissionClient struct {
	Client *permission.PermissionGRPCClient
}

func NewPermissionClient(url string) (*PermissionClient, error) {
	client, err := permission.NewPermissionGRPCClient(url)
	if err != nil {
		return nil, err
	}
	return &PermissionClient{
		Client: client,
	}, nil
}

func (c *PermissionClient) Close() {
	c.Client.Close()
}
