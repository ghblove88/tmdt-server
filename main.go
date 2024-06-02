package main

import (
	"EcdsServer/api"
	"EcdsServer/common"
	"EcdsServer/runtime"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net"
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

func startListening() {

	// 创建通道用于在服务器线程和主线程之间传递数据包
	packetChan := make(chan common.Packet)

	// 启动Socket服务器线程
	go common.RunServer(packetChan, viper.GetString("socket_server.address"),
		viper.GetInt("socket_server.port"))

	// 主线程处理数据包
	for packet := range packetChan {
		handlePacket(packet.Conn, packet)
	}
}

func main() {

	// 定义一个名为 test 的命令行参数，默认值为 false，带有描述信息
	common.TestMode = flag.Bool("test", false, "Enable test mode to run specific blocks of code")

	// 解析命令行参数
	flag.Parse()

	//启动内核服务
	rt := runtime.Runtime{}
	rt.Run()

	//启动socket 监听
	go startListening()

	//启动restful 服务
	api.Router()
}

func handlePacket(conn net.Conn, packet common.Packet) {
	fmt.Println("Received packet:")
	fmt.Println("Size:", packet.Size)
	fmt.Println("Command:", packet.Command)
	fmt.Println("Content:", packet.Content)

	// 在主线程中使用conn发送数据
	common.SendPacket(conn, "0002", "ok")
}
