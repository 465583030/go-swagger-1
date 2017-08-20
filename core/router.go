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
	*Swagger
	params  []*Param
	body    interface{}
	data    interface{}
	summary string
	desc    string
}

func (this *SwagRouter) AddPath(basePath, route, ms string) {
	this.Swagger.AddPath(basePath, route, ms, this.summary, this.desc, this.params, this.body, this.data)
}

func (this *SwagRouter) Clear() {
	this.params = make([]*Param, 0)
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
	this.desc = strings.Join(info[1:], "<br/>\n")
}

func (this *SwagRouter) QueryParam(name, desc string) *Param {
	param := &Param{"query", name, desc, "string", false, "", false, nil}
	this.params = append(this.params, param)
	return param
}

func (this *SwagRouter) PathParam(name, desc string) *Param {
	param := &Param{"path", name, desc, "string", true, "", false, nil}
	this.params = append(this.params, param)
	return param
}

func (this *SwagRouter) FileParam(name, desc string) *Param {
	param := &Param{"formData", name, desc, "file", false, "form", true, nil}
	this.params = append(this.params, param)
	return param
}

func NewSwagRouter() *SwagRouter {
	swag := new(SwagRouter)
	swag.params = make([]*Param, 0)
	swag.Swagger = NewSwagger()
	return swag
}
