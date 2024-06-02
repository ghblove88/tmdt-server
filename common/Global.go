package common

import (
	"EcdsServer/models"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	TestMode  *bool    //测试模式
	GDbsname  []string //记录系统中 洗消数据库 别名
	GDbserver string   //记录服务数据的数据库

	Config          *viper.Viper //全局配置变量；
	MsgQueueInfo    StringQueue  //短信队列
	MsgQueueWarning StringQueue  //警告短信队列
	MsgQueueError   StringQueue  //错误短信队列
)

func init() {
	//初始化日志记录
	GetLogger()
	//初始化配置文件
	viper.SetConfigFile("./app.yaml") // 指定配置文件路径
	viper.SetConfigName("app")        // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")       // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")          // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()       // 查找并读取配置文件
	if err != nil {                   // 处理读取配置文件的错误
		zap.S().Panic("Fatal error config file: ", zap.Error(err))
	}

	Config = viper.GetViper()
	models.GDRSDns = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		Config.GetString("ecds_server.username"),
		Config.GetString("ecds_server.password"), Config.GetString("ecds_server.address"),
		Config.GetInt32("ecds_server.port"), Config.GetString("ecds_server.name"))
	zap.S().Infoln("主数据库地址:", "dns:", models.GDRSDns)

	MsgQueueInfo = *NewStringQueue(50)
	MsgQueueWarning = *NewStringQueue(50)
	MsgQueueError = *NewStringQueue(50)
}
