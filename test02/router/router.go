package router

import (
	"go-web-demo/test02/handler"
	"go-web-demo/test02/middleware" // 引入中间件
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 聊天接口不需要 JWT 鉴权，先让它裸奔跑通
	r.POST("/chat", handler.ChatHandler)

	memoGroup := r.Group("/memos")
	memoGroup.Use(middleware.JWTAuth()) 
	{
		memoGroup.GET("", handler.GetMemosHandler)
		memoGroup.POST("", handler.CreateMemoHandler)
		memoGroup.PUT("/:id", handler.UpdateMemoHandler)
		memoGroup.DELETE("/:id", handler.DeleteMemoHandler)
	}

	return r
}