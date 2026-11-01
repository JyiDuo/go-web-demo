package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
//签名密钥
var jwtSecret = []byte("my-secret-key-123456")
//生成Token
func GenerateToken(userID uint) (string, error) {
	// 1. 设置载荷（Payload），也就是 Token 里存什么信息
	claims:=jwt.MapClaims{
		"user_id":userID,
		"exp":time.Now().Add(time.Hour*24).Unix(), // 过期时间，24小时后
	}
	// 2. 创建一个新的 Token 对象，指定签名方法和载荷
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	// 3. 使用指定的签名密钥生成 Token 字符串
	return token.SignedString(jwtSecret)
}
//解析Token
func ParseToken(tokenString string) (uint, error) {
	//1.解析token并且验证签名
	token,err:=jwt.Parse(tokenString,func(token*jwt.Token)(interface{},error){
		return jwtSecret,nil
	})
	if err!=nil||!token.Valid{
		return 0,errors.New("无效的Token")
	}
	//2.从token里拿到载荷
	if claims,ok:=token.Claims.(jwt.MapClaims);ok{
		userID:=uint(claims["user_id"].(float64))
		return userID,nil
	}
	return 0,errors.New("无法解析Token")
}