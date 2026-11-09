package handler

import (
	"context"
	"encoding/json" // 新增：用于解析 AI 传回来的 JSON 参数
	"fmt"
	"net/http"
	"go-web-demo/test02/model" // 新增：我们要调用数据库层了
	"go-web-demo/test02/tools"
	"go-web-demo/test02/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type ChatReq struct {
	Prompt string `json:"prompt"`
}

func ChatHandler(c *gin.Context) {
	var req ChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	ctx := context.Background()
	messages := []*schema.Message{
		{Role: schema.User, Content: req.Prompt},
	}

	// 1. 绑定工具
	toolsList := tools.GetMemoTools()
	ChatModelWithTools, err := utils.ChatModel.WithTools(toolsList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定工具失败"})
		return
	}

	// 2. 第一次调用大模型
	resp, err := ChatModelWithTools.Generate(ctx, messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 思考失败: " + err.Error()})
		return
	}

	// 3. 检查 AI 是不是想调用工具
	if len(resp.ToolCalls) > 0 {
		toolCall := resp.ToolCalls[0]
		
		if toolCall.Function.Name == "create_memo" {
			// 4. 极其关键：AI 传回来的是 JSON 字符串，我们必须解析它！
			var params tools.CreateMemoParams
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "解析 AI 参数失败"})
				return
			}
			
			// 真正存数据库（用解析出来的 params.Content）
			memo := model.CreateMemo(params.Content)

			// 5. 把工具执行的结果，追加到消息上下文里
			messages = append(messages, resp) // 先把 AI 决定调用工具的那条消息加进去
			messages = append(messages, &schema.Message{
				Role:       schema.Tool,
				Content:    fmt.Sprintf("操作成功,备忘录ID为 %d", memo.ID),
				ToolCallID: toolCall.ID,
			})

			// 6. 第二次调用大模型，让 AI 总结结果
			finalResp, err := ChatModelWithTools.Generate(ctx, messages)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 总结失败: " + err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "操作成功",
				"answer":  finalResp.Content,
			})
			return
		}
	}

	// 如果 AI 没有调用工具（比如用户只是问“你是谁”），直接返回回答
	c.JSON(http.StatusOK, gin.H{
		"message": "请求成功",
		"answer":  resp.Content,
	})
}