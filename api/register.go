package api

import (
	interstruct "ehmanager/module/datatypes"
	"ehmanager/module/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

var DB *db.MySQLDB
var SecretKey = ""

const (
	OK = http.StatusOK

	MovedPermanently = http.StatusMovedPermanently

	BadRequest   = http.StatusBadRequest
	Unauthorized = http.StatusUnauthorized
	Forbidden    = http.StatusForbidden
	NotFound     = http.StatusNotFound

	ServerError        = http.StatusInternalServerError
	ServiceUnavailable = http.StatusServiceUnavailable
)

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

const (
	TokenInvalid            CustomError = "Token 无效"
	TokenVerificationFailed CustomError = "Token 验证失败"
	TokenExpiration         CustomError = "Token 过期"
	TokenLogout             CustomError = "Token 用户退出"
)

type RouterSetupFunc func(router *gin.Engine)

// 放注册的接口
var registeredSetups []RouterSetupFunc

// RegisterSetup 注册
func RegisterSetup(setup RouterSetupFunc) {
	registeredSetups = append(registeredSetups, setup)
}

func InitRouters(engine *gin.Engine) {
	for _, setup := range registeredSetups {
		setup(engine)
	}
}

// 账号的验证
func SetDB(sqldb *db.MySQLDB) {
	DB = sqldb
}

// 异常处理
func HandleException(c *gin.Context) {
	if r := recover(); r != nil {
		log.Printf("接口错误：", r)
		c.JSON(ServerError, gin.H{"message": "接口错误", "info": r})
	}
}
func ValidateAccount(c *gin.Context) (*interstruct.LoginData, error) {
	//token转换
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			c.JSON(ServerError, ErrorResponse())
			logrus.Errorln(r)
		}
	}()
	var logininfo interstruct.LoginData

	tokenString, err := c.Cookie("token")

	if err != nil {
		//获取token失败
		logrus.Errorln("token获取失败：", err)
		//c.JSON(ServerError, ErrorResponse())
		c.Redirect(MovedPermanently, "/login")
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		//解析token失败
		logrus.Errorln("token解析失败：", err)
		//c.JSON(ServerError, ErrorResponse())
		c.Redirect(MovedPermanently, "/login")
		return nil, err
	}

	if !token.Valid {
		//token无效
		logrus.Errorln(tokenString + "--token无效")
		//c.JSON(Forbidden, ErrorResponse())
		c.Redirect(MovedPermanently, "/login")
		return nil, TokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		//token验证失败
		logrus.Errorln(tokenString + "--token无效")
		//c.JSON(Forbidden, ErrorResponse())
		c.Redirect(MovedPermanently, "/login")
		return nil, TokenVerificationFailed
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	logout := int64(claims["logout"].(float64))
	if logout != 0 {
		logrus.Println("用户退出登录")
		c.Redirect(MovedPermanently, "/login")
		return nil, TokenLogout
	}

	if time.Now().After(exp) {
		//token过期
		logrus.Infoln(tokenString + "--token过期")
		//c.JSON(Unauthorized, ErrorResponse())
		c.Redirect(MovedPermanently, "/login")
		return nil, TokenExpiration
	}

	logininfo.Username = claims["username"].(string)
	logininfo.Password = claims["password"].(string)

	return &logininfo, nil

}

func ErrorResponse() gin.H {
	return gin.H{"message": "error"}
}

func HeaderNoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	c.Header("Pragma", "no-cache")
}
