package tools

import (
	"context"
	"fmt"
	"go-web-demo/test02/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/components/tool/utils" // 注意引入了这个简化工具包
)

// 1. 定义参数结构体（保留 jsonschema 注释）
type CreateMemoParams struct {
	Content string `json:"content" jsonschema:"description=备忘录的具体内容"`
}

// 2. 实现真正干活的函数
func CreateMemoFunc(ctx context.Context, params *CreateMemoParams) (string, error) {
	// 调用你以前写好的 model 层存数据库
	memo := model.CreateMemo(params.Content)
	// 返回给 AI 的结果
	return fmt.Sprintf("创建成功,备忘录ID为:%d", memo.ID), nil
}

// 3. 注册工具（使用 utils.InferOption 自动推导）
func GetMemoTools() []*schema.ToolInfo {
    memoTool, err := utils.InferTool("create_memo", "当用户想要新建一条备忘录或记事时使用此工具", CreateMemoFunc)
    if err != nil {
        panic("工具注册失败: " + err.Error())
    }

    // 将工具转成 ToolInfo 返回
    toolInfo, err := memoTool.Info(context.Background())
    if err != nil {
        panic("获取工具信息失败: " + err.Error())
    }
    return []*schema.ToolInfo{toolInfo}
}