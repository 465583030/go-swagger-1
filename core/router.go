package core

import (
	"strings"
)

type ISwagRouterBase interface {
	QueryParam(name, desc string) *Param
	PathParam(name, desc string) *Param
	FileParam(name, desc string) *Param
}

type SwagRouter struct {
	engine  *Swagger
	params  []*Param
	body    interface{}
	data    interface{}
	summary string
	desc    string
}

func (this *SwagRouter) SetEngine(engine *Swagger) {
	this.engine = engine
}

func (this *SwagRouter) AddPath(basePath, route, ms string) {
	this.engine.AddPath(basePath, route, ms, this.summary, this.desc, this.params, this.body, this.data)
}

func (this *SwagRouter) Clear() {
	this.params = nil
	this.body = nil
	this.data = nil
	this.summary = ""
	this.desc = ""
}

func (this *SwagRouter) Params(ps ...*Param) {
	this.params = ps
}

func (this *SwagRouter) Body(body interface{}) {
	this.body = body
}

func (this *SwagRouter) Data(data interface{}) {
	this.data = data
}

func (this *SwagRouter) Info(info ...string) {
	if len(info) < 1 {
		return
	}
	this.summary = info[0]
	this.desc = strings.Join(info[1:], "\n")
}

func (this *SwagRouter) QueryParam(name, desc string) *Param {
	return &Param{"query", name, desc, "string", false, "", false, nil}
}

func (this *SwagRouter) PathParam(name, desc string) *Param {
	return &Param{"path", name, desc, "string", true, "", false, nil}
}

func (this *SwagRouter) FileParam(name, desc string) *Param {
	return &Param{"formData", name, desc, "file", false, "form", true, nil}
}
