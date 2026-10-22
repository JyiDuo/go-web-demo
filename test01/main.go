package main

import "github.com/gin-gonic/gin"

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		var req RegisterReq

		// 尝试把前端发来的 JSON 绑定到 req 结构体上
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "数据格式不对"})
			return
		}

		c.JSON(200, gin.H{
			"message": "注册成功",
			"用户":      req.Username,
			"密码":      req.Password, // 真实业务里绝不能返回密码，这里只是为了演示
		})
	})
	r.Run(":8080")
}

/*
func main(){
	r:=gin.Default()
	r.GET("/:user/:id",func(c*gin.Context){
		name:=c.Param("user")
		id:=c.Param("id")
		if name==""{
			name="陌生人"
		}
		c.JSON(200,gin.H{
			"message":"正在查询用户id",
			"查询成功":name+id,
		})
	})
	r.Run(":8080")
}
func main(){
	r:=gin.Default()
	r.GET("/hello",func(c *gin.Context){
		name:=c.Query("name")
		if name==""{
			name="陌生人"
		}
		c.JSON(200,gin.H{
			"message":"你好"+name,
		})
	})
	r.Run(":8080")
}

func main(){
	r:=gin.Default()
	r.GET("/hello",func(c *gin.Context){
		c.JSON(200,gin.H{
			"message":"这是我的第一个接口",
			"age":20,
			"is_vip":true,
			"skill":[]string{"Go","Agent"},
		})
	})
	r.Run(":8080")
}*/
