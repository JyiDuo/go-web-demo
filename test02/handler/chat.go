package handler
import (
	"context"
	"net/http"
	"go-web-demo/test02/utils"
	"github.com/cloudwego/eino/schema" // 引入 Eino 的消息结构
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
	//构造上下文并且设置超时时间
	ctx:=context.Background()
	//构造消息列表
	//模拟一下“记忆”：假设用户之前说过“我叫张三”，现在问“我叫什么？”
	// 在真实业务里，你需要把数据库里的历史聊天记录查出来，转成 []*schema.Message
	messages:=[]*schema.Message{
		{Role:schema.User,Content:"我叫张三"},
		{Role:schema.Assistant,Content:"好的，我记住了，你叫张三"},
		{Role:schema.User,Content:req.Prompt},//当前用户的问题
	}
	//调用Eino模型进行问答
	resp,err:=utils.ChatModel.Generate(ctx,messages)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"AI 思考失败:"+err.Error(),
		})
		return
	}
	//返回回答给前端
	c.JSON(http.StatusOK,gin.H{
		"message":"请求成功",
		"answer":resp.Content,
	})
}