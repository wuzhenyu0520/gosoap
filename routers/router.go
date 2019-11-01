package routers

import (
	"easemob-gosoap/controllers"
	"github.com/astaxie/beego"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    beego.Router("/*", &controllers.ProxyApiController{}, "*:AllMethod")
}
