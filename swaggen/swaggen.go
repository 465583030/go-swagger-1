package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/inu1255/go-swagger/core"
)

var (
	tmplDir    = flag.String("d", "tmpl", "template director")
	entityDir  = flag.String("e", "entity", "entity director")
	serviceDir = flag.String("s", "service", "service director")
	ext        = flag.String("ext", "go", "export file ext")
)

type Entity struct {
}

func Generate(url string) error {
	b, err := ioutil.ReadFile(url)
	if err != nil {
		return err
	}
	swag := &core.Swagger{}
	json.Unmarshal(b, swag)
	for name, entity := range swag.Definitions {
		fmt.Println(name)
		GenerateEntity(name, entity)
	}
	return nil
}

func GenerateEntity(name string, entity *core.Definition) {
	tpl, err := template.ParseFiles(*tmplDir + "/entity.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
	wr, err := os.OpenFile(*entityDir+"/"+name+"."+*ext, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Println(err)
		return
	}
	defer wr.Close()
	err = tpl.Execute(wr, entity)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	os.Mkdir(*entityDir, 0775)
	err := Generate("../api/swagger.json")
	if err != nil {
		log.Println(err)
	}
}
