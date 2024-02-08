package account

import jwtgo "github.com/dgrijalva/jwt-go"

type (
	AuthUserData struct { //common data in redis and mysql
		Uid      string `json:"uid"`
		Password string `json:"password"`
		Type     int    `json:"type"`
		Status   int    `json:"status"`
		Name     string `json:"name"`
		Active   int    `json:"active"`
		Nid      string `json:"nid"`
		Level    int    `json:"level"`
		Gid      string `json:"gid"`
	}
	JwtUserClaims struct {
		jwtgo.StandardClaims
		Info AuthUserData `json:"info"`
	}
)
