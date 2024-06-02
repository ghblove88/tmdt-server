package main

import (
	"TmdtServer/api"
	"TmdtServer/common"
	"TmdtServer/runtime"
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	// 取当前程序运行路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	zap.S().Infoln("当前程序路径:", dir)
}

func main() {

	// 定义一个名为 test 的命令行参数，默认值为 false，带有描述信息
	common.TestMode = flag.Bool("test", false, "Enable test mode to run specific blocks of code")

	// 解析命令行参数
	flag.Parse()

	//启动内核服务
	rt := runtime.Runtime{}
	rt.Run()

	//启动restful 服务
	api.Router()
}
