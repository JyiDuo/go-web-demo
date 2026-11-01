package middleware

import (
	"net/http"
	"strings"
	"go-web-demo/test02/utils"
	"github.com/gin-gonic/gin"
)
//JWT中间件，用于验证请求中的Token
func JWTAuth()gin.HandlerFunc{
	return func(c*gin.Context){
// 1. 从请求头里拿 Token。前端通常会放在 Authorization: Bearer <token> 里
		authHeader:=c.GetHeader("Authorization")
		if authHeader==""||!strings.HasPrefix(authHeader,"Bearer "){
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"请先登录(未携带Token)",
			})
			c.Abort() //阻止后续处理
			return
		}
		//2.去掉“Bearer ”前缀，拿到真正的 Token 字符串
		tokenStr:=strings.TrimPrefix(authHeader,"Bearer ")
		//3.解析 Token，拿到 userID
		userID,err:=utils.ParseToken(tokenStr)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"Token无效或已过期",
			})
			c.Abort()
			return
		}
		//4.把 userID 存到上下文里，后续的处理函数可以用 c.Get("userID") 拿到
		c.Set("user_id",userID)
		c.Next()
	}
}