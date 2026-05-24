package main

import (
	"dashboard-api/config"
	"dashboard-api/controller"
	_ "dashboard-api/db"
	"dashboard-api/middle"
	"dashboard-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化gin 对象
	r := gin.Default()
	//初始化k8s client
	service.K8s.Init()
	//初始化数据库
	// db.Init()
	r.Use(middle.Cors())
	//jwt token验证
	//r.Use(middle.JWTAuth())
	//初始化路由规则
	controller.Router.InitApiRouter(r)
	//终端websocket
	go func()  {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
	}()
	//gin程序启动
	r.Run(config.ListenAddr)
}