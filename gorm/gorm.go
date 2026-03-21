package gorm

import (
	"context"
	
	"gorm.io/gorm"
)

func New(params Params) App {
	app := &app{
		nodes: make(map[string]DBNode),
		roles: make(map[string][]DBNode),
	}
	
	for _, node := range params.DBNodes {
		app.nodes[node.Name] = node
		app.roles[node.Role] = append(app.roles[node.Role], node)
	}
	
	return app
}

type Init func(params Params) App

func (app *app) Get(context.Context) *gorm.DB {
	nodes := app.roles["primary"]
	if len(nodes) == 0 {
		return nil
	}
	
	return nodes[0].DB
}

func (app *app) GetByRole(ctx context.Context, role string) *gorm.DB {
	nodes := app.roles[role]
	if len(nodes) == 0 {
		return nil
	}
	
	return app.resolver.Resolve(nodes)
}

func (app *app) GetByName(ctx context.Context, name string) *gorm.DB {
	node, ok := app.nodes[name]
	if !ok {
		return nil
	}
	
	return node.DB
}
