package main

import (
	"go-web-demo/test02/router"
)

func main() {
	//初始化并获取路由
	r := router.SetRouter()
	//启动服务
	r.Run(":8080")
}
