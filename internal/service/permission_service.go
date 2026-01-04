package service

import (
	"context"
	"sync"

	"github.com/casbin/casbin/v2"
)

type PermissionClient interface {
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
	AddPolicy(ctx context.Context, userID, resource, action string) (bool, error)
	RemovePolicy(ctx context.Context, userID, resource, action string) (bool, error)
	AddPolicyIfNotExists(ctx context.Context, userID, resource, action string) (bool, error)
	Close() error
}

type PermissionService struct {
	enforcer *casbin.Enforcer
	mu       sync.RWMutex
}

func NewPermissionService(e *casbin.Enforcer) *PermissionService {
	return &PermissionService{
		enforcer: e,
	}
}

func (p *PermissionService) AddPolicyIfNotExists(sub, obj, act string) (bool, error) {
	hasPolicy, err := p.enforcer.HasPolicy(sub, obj, act)
	if err != nil {
		return false, err
	}

	if !hasPolicy {
		success, err := p.AddPolicy(sub, obj, act)
		return success, err
	}
	return true, nil
}

func (p *PermissionService) CheckPermission(subj, obj, act string) (bool, error) {
	p.mu.RLock()
	defer p.mu.Unlock()
	return p.enforcer.Enforce(subj, obj, act)
}

func (p *PermissionService) AddPolicy(sub, obj, act string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.enforcer.AddPolicy(sub, obj, act)
}

func (p *PermissionService) RemovePolicy(sub, obj, act string) (bool, error) {
	return p.enforcer.RemovePolicy(sub, obj, act)
}
