package router

import (
	"go-web-demo/test02/handler"
	"go-web-demo/test02/middleware" // 引入中间件
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	memoGroup := r.Group("/memos")
	// 挂载中间件：只有通过 JWTAuth 安检，才能进入下面的路由
	memoGroup.Use(middleware.JWTAuth()) 
	{
		memoGroup.GET("", handler.GetMemosHandler)
		memoGroup.POST("", handler.CreateMemoandler)
		memoGroup.PUT("/:id", handler.UpdateMemoHandler)
		memoGroup.DELETE("/:id", handler.DeleteMemoHandler)
	}

	return r
}