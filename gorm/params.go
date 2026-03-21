package gorm

import (
	"context"
	
	"gorm.io/gorm"
)

type DBNode struct {
	DB   *gorm.DB
	Name string
	Role string
}

type Params struct {
	DBNodes []DBNode
}

type App interface {
	Get(ctx context.Context) *gorm.DB
	GetByRole(ctx context.Context, role string) *gorm.DB
	GetByName(ctx context.Context, name string) *gorm.DB
}

type Resolver interface {
	Resolve(nodes []DBNode) *gorm.DB
}

type app struct {
	nodes    map[string]DBNode
	roles    map[string][]DBNode
	resolver Resolver
}
