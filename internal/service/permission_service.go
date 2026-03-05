package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/redis/go-redis/v9"
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
	watcher  persist.Watcher
	mu       sync.RWMutex
}

func NewPermissionService(e *casbin.Enforcer, cfg config.Redis) (*PermissionService, error) {
	ps := &PermissionService{
		enforcer: e,
	}
	w, err := rediswatcher.NewWatcher(cfg.Addr, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: cfg.Password,
		},
		Channel: "/casbin",
	})
	if err != nil {
		return nil, fmt.Errorf("error create csbin watcher: %w", err)
	}
	if err := e.SetWatcher(w); err != nil {
		return nil, fmt.Errorf("error setting watcher for casbin enforcer: %w", err)
	}
	w.SetUpdateCallback(func(msg string) {
		ps.mu.Lock()
		ps.enforcer.LoadPolicy()
		ps.mu.Unlock()
	})
	return ps, nil
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
	defer p.mu.RUnlock()
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

func (p *PermissionService) Close() error {
	if p.watcher != nil {
		p.watcher.Close()
	}
	return nil
}
