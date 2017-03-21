package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/inu1255/go-swagger/core"
)

var (
	tmplDir    = flag.String("d", "tmpl", "template director")
	entityDir  = flag.String("e", "", "entity director")
	serviceDir = flag.String("s", "service", "service director")
	ext        = flag.String("ext", "js", "export file ext")
	url        = flag.String("u", "", "the path or url for swagger.json")
	swag       *core.Swagger
	tpl        *template.Template
)

func Generate(url string) (err error) {
	if url == "" {
		return errors.New("need -u")
	}
	var b []byte
	if strings.HasPrefix(url, "http") {
		req, err := http.Get(url)
		if err != nil {
			return err
		}
		b, err = ioutil.ReadAll(req.Body)
	} else {
		b, err = ioutil.ReadFile(url)
	}
	if err != nil {
		return err
	}
	log.Println("load swagger.json from ", url)
	swag = &core.Swagger{}
	json.Unmarshal(b, swag)
	if *entityDir != "" {
		for name, item := range swag.Definitions {
			GenerateEntity(NewEntity(name, item))
		}
	}
	for path, table := range swag.Paths {
		for typ, item := range table {
			AddMethod(path, typ, item)
		}
	}
	for _, item := range services {
		GenerateService(item)
	}
	return nil
}

func GenerateEntity(data *Entity) {
	name := data.HungarianName()
	wr, err := os.OpenFile(*entityDir+"/"+name+"."+*ext, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Println(err)
		return
	}
	defer wr.Close()
	err = tpl.ExecuteTemplate(wr, "entity.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

func GenerateService(data *Service) {
	name := data.HungarianName()
	wr, err := os.OpenFile(*serviceDir+"/"+name+"."+*ext, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Println(err)
		return
	}
	defer wr.Close()
	err = tpl.ExecuteTemplate(wr, "service.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	os.Mkdir(*entityDir, 0775)
	os.Mkdir(*serviceDir, 0775)
	var err error
	tpl, err = template.ParseFiles(*tmplDir+"/entity.tmpl", *tmplDir+"/service.tmpl")
	if err != nil {
		log.Println(err)
	}
	err = Generate(*url)
	if err != nil {
		log.Println(err)
	}
}
