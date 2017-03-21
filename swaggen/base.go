package main

import (
	"regexp"
	"strings"
)

var (
	camel_reg     = regexp.MustCompile(`(^|\.)\w`)
	hungarian_reg = regexp.MustCompile(`\.?[A-Z]`)
)

type Base struct {
	Name string `json:"name,omitempty" xorm:"" gev:"实例名"`
}

func (this *Base) GetName() string {
	if len(this.Name) > 1 && this.Name[:1] == "." {
		return this.Name[1:]
	}
	return this.Name
}
func (this *Base) CamelName() string {
	return camel_reg.ReplaceAllStringFunc(this.GetName(), func(src string) string {
		n := len(src)
		if n > 1 {
			return strings.ToUpper(src[1:])
		}
		return strings.ToUpper(src)
	})
}
func (this *Base) HungarianName() string {
	return hungarian_reg.ReplaceAllStringFunc(this.GetName(), func(src string) string {
		n := len(src)
		if n > 1 {
			return strings.ToLower(src)
		}
		return strings.Join([]string{"_", strings.ToLower(src)}, "")
	})
}
func (this *Base) Package() string {
	return *entityDir
}
