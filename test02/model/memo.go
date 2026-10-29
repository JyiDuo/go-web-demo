package model
import(
	"gorm.io/gorm"
	 "gorm.io/driver/sqlite"  
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