package model
import(
	"gorm.io/gorm"
	"github.com/glebarez/sqlite" 
)
//定义备忘录的数据结构
type Memo struct{
	gorm.Model// gorm.Model 是 GORM 内置的结构体，
	//包含了 ID(主键)、CreatedAt(创建时间)、UpdatedAt、DeletedAt
	//GORM会自动帮你生成和管理自增ID
	Content string `json:"content"`
}
//全局的数据库连接对象（大写开头，方便其他包调用）
var DB *gorm.DB
//初始化数据库连接
func InitDB(){
	var err error
	//连接SQLite数据库，如果文件不存在，它会自动在当前目录下创建一个memo.db文件
	DB, err = gorm.Open(sqlite.Open("memo.db"), &gorm.Config{})
	if err!=nil{
		panic("连接数据库失败："+err.Error())
	}
	//自动迁移：让GORM帮我们在数据库里建好对应的表结构
	DB.AutoMigrate(&Memo{})
}
//新增备忘录
func CreateMemo(content string)Memo{
	memo:=Memo{Content: content}
	//写入数据库
	DB.Create(&memo)
	return memo
}
//获取所有备忘录
func GetAllMemos()[]Memo{
	var memos []Memo
	//查询数据库里所有的备忘录记录，塞进memos切片里
	DB.Find(&memos)
	return memos
}
//修改备忘录里的数据
func UpdateMemo(id uint,content string)(Memo,error){
	var memo Memo
	if err:=DB.First(&memo,id).Error;err!=nil{
		return memo,err
	}
	//找到了将内容替换掉
	memo.Content=content
	//保存回数据库
	DB.Save(&memo)
	return memo,nil
}
//删除备忘录
func DeleteMemo(id uint)error{
	//delete方法默认是软删除（只会打个时间戳标记，数据还在硬盘里）
	//这里我们直接传ID，GORM会去删对应的记录
	result:=DB.Delete(&Memo{},id)
	if result.RowsAffected==0{
		return gorm.ErrRecordNotFound//result.RowsAffected：这次操作影响的行数。
		//说明一条没有删除，ID 不存在
	}
	return result.Error//如果删除时报错了，这里面就是装的错误信息，如果是nil说明顺利执行
}