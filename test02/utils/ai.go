package utils
import(
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)
//Deepseek请求体结构
type DeepseekReq struct{
	Model string `json:"model"`
	Messages []Message `json:"messages"`
}
type Message struct{
	Role string `json:"role"`
	Content string `json:"content"`
}
//Deepseek响应体结构（只解析需要的那部分）
type DeepSeekResp struct{
	Choices []struct{
		Message Message `json:"message"`
	} `json:"choices"`
}
//调用大模型的核心函数
func AskDeepSeek(prompt string)(string,error){
	client:=resty.New()
	//构造请求体
	reqBody:=DeepseekReq{
		Model:"deepseek-chat",//模型名称
		Messages:[]Message{
			{
				Role:"user",
				Content:prompt,
			},
		},
	}
	 apiKey := os.Getenv("DEEPSEEK_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("环境变量 DEEPSEEK_API_KEY 未配置")
	}
	var respData DeepSeekResp
	//发送post请求
	resp,err:=client.R().SetHeader("Content-Type","application/json").
SetHeader("Authorization","Bearer "+apiKey).
SetBody(reqBody).SetResult(&respData).Post("https://api.deepseek.com/chat/completions")
	if err!=nil{
		return "",err
	}
	if resp.StatusCode()!=200{
		return "",fmt.Errorf("API请求失败,状态码:%d,返回内容:%s",resp.StatusCode(),resp.String())
	}
	//返回模型回答的内容
	if len(respData.Choices)>0{
		return respData.Choices[0].Message.Content,nil

	}
	return "",fmt.Errorf("未获取到模型回答")
}
