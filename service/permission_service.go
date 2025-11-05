package service

import "github.com/casbin/casbin/v2"

type PermissionService struct {
	enforcer *casbin.Enforcer
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

func (p *PermissionService) CheckPermission(userID, resource, act string) (bool, error) {
	return p.enforcer.Enforce(userID, resource, act)
}

func (p *PermissionService) AddPolicy(policy, resource, permission string) (bool, error) {
	return p.enforcer.AddPolicy(policy, resource, permission)
}

func (p *PermissionService) DeletePolicy(policy, resource, permission string) (bool, error) {
	return p.enforcer.RemovePolicy(policy, resource, permission)
}
