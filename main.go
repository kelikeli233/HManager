package ehmanager

import (
	"ehmanager/api"
	"ehmanager/module/db"
	config "ehmanager/module/default"
	"ehmanager/module/key"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func main() {
	//日志
	logSave, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal("无法创建日志文件：", err)
	}
	defer logSave.Close()
	logrus.SetFormatter(&LogFormatter{})
	logrus.SetOutput(logSave)

	args, err := config.LoadStartConfig()
	if err != nil {
		return
	}

	if args.Other.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err != nil {
		logrus.Println(err)
		return
	}
	//
	//数据库
	dbDriver, err := db.Connect(args.Database.Backend, db.ConnectionArgs{
		MasterDSN:              args.Database.Master.DSN,
		ReplicaDSN:             args.Database.Replica.DSN,
		MaxIdleConns:           args.Database.MaxIdleConns,
		MaxOpenConns:           args.Database.MaxOpenConns,
		MaxConnLifetimeSeconds: args.Database.MaxConnLifetimeSeconds,
	})
	//数据库
	if err != nil {
		logrus.Println(err)
		os.Exit(101)
	}
	//syslog

	//syslog

	// 设置jwt密钥

	logrus.Println("设置JWT密钥")
	SecretKey, err := key.RandomKey()
	if err != nil {
		logrus.Errorln("密钥设置错误")
		return
	}
	api.SecretKey = SecretKey
	fmt.Println("设置JWT密钥为：", SecretKey)
	logrus.Println("设置JWT密钥为：", SecretKey)

	// 开启服务端
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.Static("/public", "./public")
	router.LoadHTMLGlob("public/*.html")

	router.Use(gin.LoggerWithWriter(&GinLogWriter{logSave}))
	api.SetDB(dbDriver)
	api.InitRouters(router)
	router.Run(args.HTTP.Address)
}

type LogFormatter struct{}

type GinLogWriter struct {
	out io.Writer
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatted := fmt.Sprintf("[%s][%s]:%s\n",
		entry.Time.Format(time.RFC3339),
		entry.Level.String(),
		entry.Message)

	return []byte(formatted), nil
}

func (f *GinLogWriter) Write(p []byte) (n int, err error) {
	logMessage := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02T15:04:05-07:00"), string(p))
	_, err = f.out.Write([]byte(logMessage))
	if err != nil {
		return 0, nil
	}
	return len(p), nil
}
