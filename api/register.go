package api

import (
	"ehmanager/module/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

// 初始化
func InitRouters(router *gin.Engine) {

	for _, setup := range registeredSetups {
		setup(router)
	}
}

func ErrorResponse() gin.H {
	return gin.H{"message": "error"}
}

func HeaderNoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	c.Header("Pragma", "no-cache")
}
