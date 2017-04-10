package core

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CopySwagger() {
	var _gev_path string
	_, file, _, _ := runtime.Caller(0)
	if index := strings.LastIndex(file, "/"); index > 0 {
		_gev_path = file[:index-5]
	}
	if info, err := os.Stat("api"); err != nil || !info.IsDir() {
		cmd := exec.Command("cp", "-R", _gev_path+"/api", "api")
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("swagger文件夹已经存在")
	}
}
