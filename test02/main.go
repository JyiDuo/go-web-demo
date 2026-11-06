package main

import (
	"context" // 引入 context
	"go-web-demo/test02/model"
	"go-web-demo/test02/router"
	"go-web-demo/test02/utils"
)

func main() {
	model.InitDB()
	//初始化Eino的大模型
	ctx:=context.Background()
	if err:=utils.InitEinoModel(ctx);err!=nil{
		panic("Eino 模型初始化失败:"+err.Error())
	}
	r := router.SetupRouter()
	r.Run(":8080")
}