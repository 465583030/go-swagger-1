package core

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

var (
	MapSchema = map[string]interface{}{"$ref": "#/definitions/.map"}
	reg_route = regexp.MustCompile(`:(\w+)`)
)

type Swagger struct {
	Swagger             string                   `json:"swagger,omitempty"`
	Info                Info                     `json:"info,omitempty"`
	Host                string                   `json:"host,omitempty"`
	BasePath            string                   `json:"basePath,omitempty"`
	Schemes             []string                 `json:"schemes,omitempty"`
	Consumes            []string                 `json:"consumes,omitempty"`
	Produces            []string                 `json:"produces,omitempty"`
	SecurityDefinitions map[string]*Security     `json:"securityDefinitions,omitempty"`
	Security            []map[string]interface{} `json:"security,omitempty"`
	Paths               map[string]Path          `json:"paths,omitempty"`
	Definitions         map[string]*Definition   `json:"definitions,omitempty"`
	Tags                []*Tag                   `json:"tags,omitempty"`
}

func (this *Swagger) AddPath(basePkg, route, ms, summary, desc string, params []*Param, body interface{}, data interface{}) {
	method := NewMethod(summary, desc)
	if len(basePkg) > 1 {
		method.Tags = append(method.Tags, basePkg[1:])
	}
	if params != nil {
		count := len(params)
		method.Parameters = make([]*Param, count, count+1)
		for i := 0; i < count; i++ {
			method.Parameters[i] = params[i]
		}
	}
	if body != nil {
		this.methodBody(method, body)
	}
	if data != nil {
		method.SetResponse(this.Schema(reflect.ValueOf(data)))
	}
	route = reg_route.ReplaceAllString(route, "{$1}")
	ps := strings.Replace(basePkg+route, "//", "/", -1)
	if path, ok := this.Paths[ps]; ok {
		path.SetMethod(ms, method)
	} else {
		path = make(Path)
		path.SetMethod(ms, method)
		this.Paths[ps] = path
	}
}

func (this *Swagger) AddTag(tag string, summary ...string) {
	if len(summary) > 0 {
		this.Tags = append(this.Tags, &Tag{tag, strings.Join(summary, ";")})
	}
}

func (this *Swagger) methodBody(method *Method, v interface{}) {
	schema := this.Schema(reflect.ValueOf(v))
	param := new(Param)
	param.In = "body"
	param.Name = "body"
	param.Schema = schema
	method.Parameters = append(method.Parameters, param)
}

// 获取 schema {type:"",schema:{$ref},items:{$ref}}
func (this *Swagger) Schema(v reflect.Value) map[string]interface{} {
	schema := make(map[string]interface{})
	switch v.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Float64, reflect.Float32:
		schema["type"] = "number"
	case reflect.Struct:
		if m, ok := v.Interface().(encoding.TextMarshaler); ok {
			schema["type"] = "string"
			b, _ := m.MarshalText()
			schema["default"] = string(b)
		} else {
			schema["$ref"] = this.Define(v)
		}
	case reflect.Interface:
		if !v.IsNil() {
			schema = this.Schema(v.Elem())
		}
	case reflect.Ptr:
		if v.IsNil() {
			schema = this.Schema(reflect.New(v.Type().Elem()).Elem())
		} else {
			schema = this.Schema(v.Elem())
		}
	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		if v.Len() > 0 {
			schema["items"] = this.Schema(v.Index(0))
		} else {
			schema["items"] = this.Schema(reflect.New(v.Type().Elem()).Elem())
		}
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Map:
		schema["$ref"] = this.Define(v)
	}
	return schema
}

// 获取 tag 的 字段名
func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

// 为一个 definition 定义属性
func (this *Swagger) defineProperty(define *Definition, v reflect.Value) {
	switch v.Kind() {
	case reflect.Interface:
		this.defineProperty(define, v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			this.defineProperty(define, reflect.New(v.Type().Elem()).Elem())
		} else {
			this.defineProperty(define, v.Elem())
		}
	case reflect.Struct:
		t := v.Type()
		n := t.NumField()
		for i := 0; i < n; i++ {
			field := t.Field(i)
			// 如果是父对象 获取父对象的属性
			if field.Anonymous && field.PkgPath == "" {
				this.defineProperty(define, v.Field(i))
			} else if !field.Anonymous && field.PkgPath == "" {
				jsonTag := field.Tag.Get("json")
				if jsonTag == "-" {
					continue
				}
				gevTag := field.Tag.Get("gev")
				if gevTag == "-" {
					continue
				}

				// 如果可以读取，获取 json的键名
				name := parseTag(jsonTag)
				if name == "" {
					name = field.Name
				}
				schema := this.Schema(v.Field(i))
				schema["description"] = gevTag
				define.Properties[name] = schema
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			define.Properties[key.String()] = this.Schema(v.MapIndex(key))
		}
	}
}

// 在 swagger.definitions 中定义 v
func (this *Swagger) Define(v reflect.Value) string {
	var pkg string
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return this.Define(v.Elem())
	case reflect.Struct:
		t := v.Type()
		// github.com.inu1255.youyue.common.models.User
		pkg = strings.Replace(t.PkgPath(), "/", ".", -1)
		pkg += "." + t.Name()
	case reflect.Map:
		pkg = fmt.Sprintf("#/definitions/.map.%v", v.Pointer())
	default:
		return "#/definitions/.map"
	}
	ds := this.Definitions
	if _, ok := ds[pkg]; !ok {
		define := NewDefinition()
		ds[pkg] = define
		this.defineProperty(define, v)
	}
	return "#/definitions/" + pkg
}

// 写入 json 文件
func (this *Swagger) WriteJson(filename string) error {
	b, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func NewSwagger() *Swagger {
	swagger := &Swagger{
		Swagger:  "2.0",
		Info:     Info{"1.0.0", "接口文档", "点击右上角`Authorize`切换token\n"},
		Host:     "",
		BasePath: "/",
		Schemes:  []string{"http"},
		Consumes: []string{"application/json", "multipart/form-data", "text/plain"},
		Produces: []string{"application/json"},
		SecurityDefinitions: map[string]*Security{
			"xauth": &Security{
				Type: "apiKey",
				Name: "access_token",
				In:   "query",
			},
		},
		Security:    make([]map[string]interface{}, 1),
		Paths:       make(map[string]Path),
		Definitions: make(map[string]*Definition),
		Tags:        make([]*Tag, 0, 3),
	}
	swagger.Security[0] = map[string]interface{}{
		"xauth": make([]string, 0),
	}
	swagger.Definitions[".map"] = NewDefinition()
	return swagger
}
