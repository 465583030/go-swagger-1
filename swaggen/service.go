package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/inu1255/go-swagger/core"
)

var (
	services       = make([]*Service, 0)
	path_reg       = regexp.MustCompile(`/\{\w+\}`)
	path_camel_reg = regexp.MustCompile(`/\w`)
)

type Service struct {
	Base
	Methods []*Method `json:"methods,omitempty" xorm:"" gev:"单个接口"`
}

func (this *Service) AddMethod(method *Method) {
	this.Methods = append(this.Methods, method)
}

type Method struct {
	core.Method
	Path string `json:"path,omitempty" xorm:"" gev:"接口uri"`
	Type string `json:"type,omitempty" xorm:"" gev:"类型get/post"`
	data *Entity
	body *Entity
}

func (this *Method) HungarianName() string {
	path := path_reg.ReplaceAllString(this.Path, "")
	return strings.Replace(path, "/", "_", -1)
}
func (this *Method) CamelName() string {
	path := path_reg.ReplaceAllString(this.Path, "")
	return path_camel_reg.ReplaceAllStringFunc(path, func(src string) string {
		return strings.ToUpper(src[1:])
	})
}
func (this *Method) Params() []*core.Param {
	res := make([]*core.Param, 0, len(this.Parameters))
	for _, item := range this.Parameters {
		if item.In != "body" {
			res = append(res, item)
		}
	}
	return res
}
func (this *Method) LowerType() string {
	return strings.ToLower(this.Type)
}
func (this *Method) IsPost() bool {
	return this.LowerType() == "post"
}
func (this *Method) Uri(pb, pe, qb, qe string) string {
	path := this.Path
	querys := make([]string, 0)
	for _, item := range this.Parameters {
		switch item.In {
		case "path":
			path = strings.Replace(path, "{"+item.Name+"}", pb+item.Name+pe, 1)
		case "query":
			querys = append(querys, item.Name+"="+qb+item.Name+pe)
		}
	}
	if len(querys) < 1 {
		return path
	}
	return path + "?" + strings.Join(querys, "&")
}
func (this *Method) HasData() bool {
	return this.GetData() != nil
}
func (this *Method) GetData() *Entity {
	if this.data != nil {
		return this.data
	}
	if this.Responses == nil {
		return nil
	}
	if data, ok := this.Responses["200"]; ok {
		if schema, ok := data.Schema.(map[string]interface{}); ok {
			ref := schema["$ref"]
			if ref != nil && len(ref.(string)) > 14 {
				name := ref.(string)[14:]
				item, ok := swag.Definitions[name]
				if ok {
					this.data = NewEntity(name, item)
					return this.data
				} else {
					log.Println(name, "未定义")
				}
			}
		}
	}
	return nil
}
func (this *Method) HasBody() bool {
	return this.GetBody() != nil
}
func (this *Method) GetBody() *Entity {
	if this.body != nil {
		return this.body
	}
	for _, item := range this.Parameters {
		if item.In == "body" {
			ref := item.Schema["$ref"]
			if ref != nil && len(ref.(string)) > 14 {
				name := ref.(string)[14:]
				item, ok := swag.Definitions[name]
				if ok {
					this.body = NewEntity(name, item)
					return this.body
				} else {
					log.Println(name, "未定义")
					break
				}
			}
		}
	}
	return nil
}

func AddMethod(path, typ string, item *core.Method) {
	method := new(Method)
	method.Method = *item
	method.Type = typ
	var name string
	index := strings.Index(path[1:], "/")
	if index > 0 {
		name = path[1 : index+1]
		method.Path = path[index+2:]
	} else {
		name = path[1:]
		method.Path = path
	}
	count := len(services)
	needNew := true
	for i := 0; i < count; i++ {
		if services[i].Name == name {
			services[i].AddMethod(method)
			needNew = false
			break
		}
	}
	if needNew {
		service := new(Service)
		service.Name = name
		service.AddMethod(method)
		services = append(services, service)
	}
}
