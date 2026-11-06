package utils

import (
	"context"
	"fmt"
	"os"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
)
//定义一个全局的chatmodel
var ChatModel *deepseek.ChatModel
//InitRinoModel初始化Eino模型
func InitEinoModel (ctx context.Context)error{
	apiKey:=os.Getenv("DEEPSEEK_API_KEY")
	if apiKey==""{
		return fmt.Errorf("环境变量 DEEPSEEK_API_KEY 未配置，请先设置后再启动")
	}
	var err error
	// 使用 Eino 的 DeepSeek 适配器创建模型实例
	ChatModel,err=deepseek.NewChatModel(ctx,&deepseek.ChatModelConfig{
		APIKey: apiKey,
		Model:"deepseek-chat",
		BaseURL:"https://api.deepseek.com",
	})
	if err!=nil{
		return err
	}
	return nil
}