package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	userID := uint(1)
	ctx := context.Background()

	// 1. 组装历史记忆
	messages := []*schema.Message{}
	history := model.GetRecentChats(userID, 10)
	for _, chat := range history {
		role := schema.User
		if chat.Role == "assistant" {
			role = schema.Assistant
		}
		messages = append(messages, &schema.Message{Role: role, Content: chat.Content})
	}
	messages = append(messages, &schema.Message{Role: schema.User, Content: req.Prompt})
	
	model.SaveChat(userID, "user", req.Prompt)

	// 2. 绑定工具
	toolsList := tools.GetMemoTools()
	ChatModelWithTools, err := utils.ChatModel.WithTools(toolsList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定工具失败"})
		return
	}

	// 3. 设置 SSE 响应头（告诉浏览器，我要持续吐数据了）
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 4. 第一次调用：流式请求
	// 注意这里用的是 Stream，而不是之前的 Generate
	stream, err := ChatModelWithTools.Stream(ctx, messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 思考失败: " + err.Error()})
		return
	}
	defer stream.Close() // 极其重要！防止内存泄漏

	// 5. 读取流式返回的数据
	var fullRespContent string // 用来把 AI 说的话拼接起来，最后存数据库
	var toolCalls []schema.ToolCall // 如果触发了工具，要在这里收集

	for {
		chunk, err := stream.Recv() // 每次接收一个小碎片
		if err != nil {
			if err == io.EOF { // 读完了
				break
			}
			fmt.Printf("读取流失败: %v\n", err)
			return
		}

		// 如果这个碎片里包含文本内容，直接推给前端（实现打字机效果）
		if chunk.Content != "" {
			fullRespContent += chunk.Content
			// SSE 格式要求：data: xxx\n\n
			c.Writer.WriteString("data: " + chunk.Content + "\n\n")
			c.Writer.Flush() // 强制把数据推给前端
		}

		// 如果这个碎片里包含工具调用，先收集起来，等流结束后处理
		if len(chunk.ToolCalls) > 0 {
			toolCalls = append(toolCalls, chunk.ToolCalls...)
		}
	}

	// 6. 处理工具调用（如果 AI 在流式输出中决定调用工具）
	if len(toolCalls) > 0 {
		toolCall := toolCalls[0] // 假设只调用一个工具

		if toolCall.Function.Name == "create_memo" {
			var params tools.CreateMemoParams
			json.Unmarshal([]byte(toolCall.Function.Arguments), &params)
			memo := model.CreateMemo(params.Content)

			// 注意：流式模式下，如果决定调工具，前面通常不会吐出字。
			// 我们需要把上下文补齐，再发第二次请求
			messages = append(messages, &schema.Message{
				Role:      schema.Assistant,
				Content:   "",
				ToolCalls: []schema.ToolCall{toolCall},
			})
			messages = append(messages, &schema.Message{
				Role:       schema.Tool,
				Content:    fmt.Sprintf("操作成功，备忘录ID为 %d", memo.ID),
				ToolCallID: toolCall.ID,
			})

			// 第二次流式调用，让 AI 根据工具结果生成最终回答
			stream2, _ := ChatModelWithTools.Stream(ctx, messages)
			defer stream2.Close()

			fullRespContent = "" // 清空之前的
			for {
				chunk, err := stream2.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					return
				}
				if chunk.Content != "" {
					fullRespContent += chunk.Content
					c.Writer.WriteString("data: " + chunk.Content + "\n\n")
					c.Writer.Flush()
				}
			}
		}
	}

	// 7. 最后把 AI 的回复存入数据库
	if fullRespContent != "" {
		model.SaveChat(userID, "assistant", fullRespContent)
	}

	// 发送结束标志
	c.Writer.WriteString("data: [DONE]\n\n")
	c.Writer.Flush()
}