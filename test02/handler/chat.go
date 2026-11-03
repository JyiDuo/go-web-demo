package handler
import(
	"net/http"
	"go-web-demo/test02/utils"
	"github.com/gin-gonic/gin"
)
type ChatReq struct{
	Prompt string `json:"prompt"`
}
func ChatHandler(c*gin.Context){
	var req ChatReq
	//接收前端传来的问题
	if err :=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"请求参数错误",
		})
		return
	}
	//调用大模型接口获取回答
	answer,err:=utils.AskDeepSeek(req.Prompt)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"调用大模型接口失败:"+err.Error(),
		})
		return
	}
	//返回回答给前端
	c.JSON(http.StatusOK,gin.H{
		"message":"请求成功",
		"answer":answer,
	})
}