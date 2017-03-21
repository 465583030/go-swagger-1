package main

import "github.com/inu1255/go-swagger/core"

type Entity struct {
	core.Definition
	Base
}

func NewEntity(name string, item *core.Definition) *Entity {
	entity := new(Entity)
	entity.Name = name
	entity.Definition = *item
	return entity
}
