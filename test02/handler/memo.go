package handler
import(
	"net/http"
	"go-web-demo/test02/model"
	"github.com/gin-gonic/gin"
)
//接收前端发来的JSON的模具
type CreateMemoReq struct{
	Content string `json:"content"`
}
//处理POST请求：新增备忘录
func CreateMemoandler(c*gin.Context){
	var req CreateMemoReq
	//1.绑定前端传来的JSON
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"内容格式不对",
		})
		return
	}
	//2.调用model层的业务逻辑
	memo:=model.CreateMemo(req.Content)
	//3.返回成功信息
	c.JSON(http.StatusOK,gin.H{
		"message":"创建成功",
		"data":memo,
	})
}
//处理GET请求：获取所有备忘录
func GetMemosHandler(c*gin.Context){
	//直接调用model层拿数据
	memos:=model.GetAllMemos()
	c.JSON(http.StatusOK,gin.H{
		"message":"获取成功",
		"data":memos,
	})
}