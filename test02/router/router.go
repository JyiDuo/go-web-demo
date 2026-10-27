package router
import(
	"go-web-demo/test02/handler"
	"github.com/gin-gonic/gin"
)
//初始化路由的方法
func SetRouter()*gin.Engine{
	r:=gin.Default()
	//把路由分组。/memos 下面的都归这里管
	memoGroup:=r.Group("/memos")
	{
		//获取列表
		memoGroup.GET("",handler.GetMemosHandler)
		//创建备忘录
		memoGroup.POST("",handler.CreateMemoandler)
	}
	return r
}