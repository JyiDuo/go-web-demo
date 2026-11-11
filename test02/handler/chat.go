package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"go-web-demo/test02/model"
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

	// 假设当前登录用户 ID 为 1
	userID := uint(1)

	ctx := context.Background()

	// ================= 1. 记忆注入开始 =================
	messages := []*schema.Message{}

	// 去数据库查最近的 10 条记录
	history := model.GetRecentChats(userID, 10)
	for _, chat := range history {
		// 把数据库里的记录，转换成大模型认识的 Message
		role := schema.User
		if chat.Role == "assistant" {
			role = schema.Assistant
		}
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: chat.Content,
		})
	}
	
	// 把当前用户问的话，追加到历史记录后面
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: req.Prompt,
	})
	// ================= 记忆注入结束 =================

	// 保存用户当前这句话到数据库
	model.SaveChat(userID, "user", req.Prompt)

	// 2. 绑定工具
	toolsList := tools.GetMemoTools()
	ChatModelWithTools, err := utils.ChatModel.WithTools(toolsList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定工具失败"})
		return
	}

	// 3. 第一次调用大模型
	resp, err := ChatModelWithTools.Generate(ctx, messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 思考失败: " + err.Error()})
		return
	}

	// 4. 检查 AI 是不是想调用工具
	if len(resp.ToolCalls) > 0 {
		toolCall := resp.ToolCalls[0]

		if toolCall.Function.Name == "create_memo" {
			var params tools.CreateMemoParams
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "解析 AI 参数失败"})
				return
			}

			memo := model.CreateMemo(params.Content)

			messages = append(messages, resp) // 先把 AI 决定调工具的消息加进去
			messages = append(messages, &schema.Message{
				Role:       schema.Tool,
				Content:    fmt.Sprintf("操作成功，备忘录ID为 %d", memo.ID),
				ToolCallID: toolCall.ID,
			})

			// 第二次调用大模型，让 AI 总结结果
			finalResp, err := ChatModelWithTools.Generate(ctx, messages)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 总结失败: " + err.Error()})
				return
			}

			// 保存 AI 的最终总结到数据库
			model.SaveChat(userID, "assistant", finalResp.Content)

			c.JSON(http.StatusOK, gin.H{
				"message": "操作成功",
				"answer":  finalResp.Content,
			})
			return
		}
	}

	// 如果 AI 没有调用工具，直接回答
	// 保存 AI 的普通回答到数据库
	model.SaveChat(userID, "assistant", resp.Content)

	c.JSON(http.StatusOK, gin.H{
		"message": "请求成功",
		"answer":  resp.Content,
	})
}