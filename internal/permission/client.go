package permission

import "context"

type PermissionClient interface {
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
	AddPolicy(ctx context.Context, userID, resource, action string) (bool, error)
	RemovePolicy(ctx context.Context, userID, resource, action string) (bool, error)
	AddPolicyIfNotExists(ctx context.Context, userID, resource, action string) (bool, error)
	Close() error
}
