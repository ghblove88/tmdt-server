//go:build YMC60读卡器
// +build YMC60读卡器

package runtime

import (
	"EcdsServer/common"
	"container/list"
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type STRUCT_READID_MSG struct {
	IP  string
	MSG string
}

func (rm *STRUCT_READID_MSG) GetType() (nType int) {
	str := strings.Split(rm.IP, ".")
	i, _ := strconv.Atoi(str[3])
	return i
}

type STRUCT_READERID_QUEUE struct {
	sem  chan int
	list list.List
}

func (rq *STRUCT_READERID_QUEUE) Push_ReadID_Queue(msg STRUCT_READID_MSG) {
	rq.sem <- 1
	_ = rq.list.PushFront(msg)
	<-rq.sem
}

func (rq *STRUCT_READERID_QUEUE) PullBcak_ReadID_Queue() (STRUCT_READID_MSG, error) {
	if rq.list.Len() <= 0 {
		return STRUCT_READID_MSG{}, errors.New("null")
	}
	rq.sem <- 1
	str := rq.list.Back()
	rq.list.Remove(str)
	<-rq.sem
	return str.Value.(STRUCT_READID_MSG), nil
}

type Socket_Reader struct {
	Readid_queue STRUCT_READERID_QUEUE
}

func (sr *Socket_Reader) Run() {
	sr.Readid_queue = STRUCT_READERID_QUEUE{make(chan int, 1), list.List{}}
	go sr.Socket_start()
}

func (sr *Socket_Reader) Socket_start() {
	service := common.Config.GetString("reader.port")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+service)
	if err != nil {
		zap.S().Errorln("Fatal error: %s", zap.Error(err))
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		zap.S().Errorln("Fatal error: %s", zap.Error(err))
		os.Exit(1)
	}

	zap.S().Infoln("Socket_Reader Running......" + tcpAddr.String())
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClient(conn)
	}
}

func (sr *Socket_Reader) Pull_ReadID_Queue() (msg STRUCT_READID_MSG, err error) {
	return sr.Readid_queue.PullBcak_ReadID_Queue()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	log.Println("From:" + conn.RemoteAddr().String() + " connection is successful。")
	for {

		response := make([]byte, 512) // 创建足够大的buffer接收数据
		length, err := conn.Read(response)
		if err != nil {
			fmt.Println("Failed to read from connection:", err)
			break
		}
		if length <= 0 {
			continue
		}

		fmt.Println("Received data:", hex.EncodeToString(response[:length]))
		no := parseCardNumberFromPacket(response[:length])
		//msg := common.Substring("0000000000"+strconv.Itoa(int(no)), len(strconv.Itoa(int(no))), 10)

		remoteip := strings.Split(conn.RemoteAddr().String(), ":")
		G_Socket_Reader.Readid_queue.Push_ReadID_Queue(STRUCT_READID_MSG{remoteip[0], no})
		log.Println(" time: " + time.Now().String() + ", From:" + conn.RemoteAddr().String() + " Receive credit card: " + no)

	}
	log.Println("From:" + conn.RemoteAddr().String() + " disconnect。")
}

// parseCardNumberFromPacket 解析自动读卡数据包，提取卡号
func parseCardNumberFromPacket(data []byte) string {
	if len(data) < 11 { // 检查数据长度是否足够
		fmt.Println("Data packet too short")
		return ""
	}

	// 确保包头和包尾正确
	if data[0] != 0xCC || data[len(data)-2] != 0xDD {
		fmt.Println("Invalid packet start or end")
		return ""
	}

	// 检查命令码是否为0x60
	if data[7] != 0x60 {
		fmt.Println("Not a card number packet")
		return ""
	}

	// 包长度字段，根据这个长度可以计算卡号长度
	packetLength := data[6]
	if int(packetLength) != len(data)-7 { // 包长度应等于总数据长度减去起始字节、包尾和校验值
		fmt.Println("Invalid packet length")
		return ""
	}

	// 提取卡号
	cardNumberStartIndex := 9
	cardNumberEndIndex := len(data) - 2 // 排除包尾和校验值
	cardNumber := data[cardNumberStartIndex:cardNumberEndIndex]
	fmt.Println("Card Number:", hex.EncodeToString(cardNumber))
	// 十六进制转十进制字符串
	cardNumberBigInt := new(big.Int)
	cardNumberBigInt.SetString(hex.EncodeToString(cardNumber), 16)
	fmt.Println("Card Number (decimal):", cardNumberBigInt.String())

	return cardNumberBigInt.String()
}
