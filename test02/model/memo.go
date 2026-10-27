package model
//1.定义备忘录的数据结构
type Memo struct{
	ID      int    `json:"id"`//序列化成 JSON 时，这个字段的键名用 "id"
	Content string `json:"content"`
}
//2.一个数据库，用切片来存数据
var MemoList []Memo
//3.核心业务逻辑：新增备忘录
func CreateMemo(content string)Memo{
	newID:=len(MemoList)+1
	newMemo:=Memo{
		ID: newID,
		Content: content,
	}
	MemoList=append(MemoList, newMemo)
	return newMemo
}
//4.核心业务逻辑：获取所有备忘录
func GetAllMemos()[]Memo{
	return MemoList
}