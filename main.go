package main

import (
	_ "easemob-gosoap/routers"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = beego.AppConfig.String("log_path")

	// map 转 json
	configstr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}
	beego.SetLogger(logs.AdapterFile, string(configstr))
	beego.SetLogFuncCall(true)
	fmt.Println(string(configstr))
	return
}


func main() {
	initLogger()
	beego.Run()
}

