package handler
import(
	"net/http"
	"go-web-demo/test02/model"
	"github.com/gin-gonic/gin"
	"strconv"
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
type UpdateMemoReq struct{
	Content string `json:"content"`
}
//处理PUT请求：修改备忘录
func UpdateMemoHandler(c*gin.Context){
	//1.从URL路径里拿ID
	idStr:=c.Param("id")
	//把字符串ID转成无符号整数
	id,err:=strconv.ParseUint(idStr,10,32)
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"非法的备忘录 ID",
		})
		return
	}
	//2.解析前端传来的新内容
	var req UpdateMemoReq
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"内容格式不对",
		})
		return
	}
	//3.调用model层去修改
	memo,err:=model.UpdateMemo(uint(id),req.Content)
	if err!=nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error":"备忘录不存在",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"修改成功",
		"data":memo,
	})
}
//处理DELETE请求：删除备忘录
func DeleteMemoHandler(c*gin.Context){
	//1.拿ID
	idStr:=c.Param("id")
	//把字符串ID转成无符号整数
	id,err:=strconv.ParseUint(idStr,10,32)
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"非法的备忘录 ID",
		})
		return
	}
	//调用model去删
	if err:=model.DeleteMemo(uint(id));err!=nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error":"备忘录不存或已删除",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"删除成功",
	})
}