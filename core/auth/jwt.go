package common

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type Admin string

const (
	AdminKey    Admin = "adminKey"
	AdminSecret Admin = "adminSecret"
)

var JwtKey []byte

type JwtUser struct {
	Id          uint64 //用户id
	OrgId       int64  //登录的机构id
	UserAccount string //用户账户
	UserName    string //用户名称
	UserMobile  string //用户手机号
	UserCode    string //用户唯一编码
}

type claimsUser struct {
	*JwtUser
	*jwt.RegisteredClaims
}

func GenerateToken(appSecret string, user *JwtUser) (string, error) {
	JwtKey = []byte(appSecret) //为密钥赋值

	//不让token过期，过期由会话来保持
	claims := &claimsUser{
		JwtUser: user,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	//生成token编码
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JwtKey)

	return token, err
}

func ParseToken(appSecret Admin, token string) (*JwtUser, error) {
	JwtKey = []byte(appSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &claimsUser{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		log.Println("parse with claims err=", err)
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*claimsUser); ok && tokenClaims.Valid {
			return claims.JwtUser, nil //解码之后的数据
		}
	}

	return nil, err
}
