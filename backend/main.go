package main

import (
	_ "api/controllers"
	_ "api/routers"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	_ "net/http"
)

func main() {

	sqlconn, _ := beego.AppConfig.String("sqlconn")
	orm.RegisterDataBase("mydatabase", "postgres", sqlconn)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Connection", "Upgrade"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Log.Outputs["console"] = ""
	beego.Run()
}
