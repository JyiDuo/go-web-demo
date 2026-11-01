package main

import (
	"fmt"
	"go-web-demo/test02/model"
	"go-web-demo/test02/router"
	"go-web-demo/test02/utils"
)

func main() {
	model.InitDB()

	// 【临时测试代码】生成一个 Token 方便测试
	testToken, _ := utils.GenerateToken(1)
	fmt.Println("====== 测试用的 Token ======")
	fmt.Println(testToken)
	fmt.Println("============================")

	r := router.SetupRouter()
	r.Run(":8080")
}