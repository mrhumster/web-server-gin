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
		success, err := p.enforcer.AddPolicy(sub, obj, act)
		return success, err
	}
	return true, nil
}

func (p *PermissionService) CheckPermission(userID, fullResource, act string) (bool, error) {
	return p.enforcer.Enforce(userID, fullResource, act)
}

func (p *PermissionService) AddRule(policy, resource, permission string) (bool, error) {
	return p.enforcer.AddPolicy(policy, resource, permission)
}

func (p *PermissionService) DeleteRule(policy, resource, permission string) (bool, error) {
	return p.enforcer.RemovePolicy(policy, resource, permission)
}
