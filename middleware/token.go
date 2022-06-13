package middleware

import (
	"bytedance_SimpleDouyin/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Key = []byte("byte dance 11111 return")

type DyClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// CreateToken 生成token
func CreateToken(userId uint, userName string) (string, error) {

	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	NowTime := time.Now()                        //当前时间
	cliaims := DyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  NowTime.Unix(),    //当前时间戳
			Issuer:    "henrik",          //颁发者签名
			Subject:   "userToken",       //主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, cliaims)
	return tokenStruct.SignedString(Key)
}

// CheckToken 检查token
func CheckToken(token string) (*DyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &DyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if Key, _ := tokenObj.Claims.(*DyClaims); tokenObj.Valid {
		return Key, true
	} else {
		return nil, false
	}

}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, common.Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := CheckToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort() //阻止执行
			return
		}
		c.Set("username", tokenStruck.UserName)
		c.Set("user_id", tokenStruck.UserId)

		c.Next()
	}
}
