//go:build CX522读卡器
// +build CX522读卡器

// Package ecd_runtime go:YMC60 老崔定制读卡器
package runtime

import (
	"EcdsServer/common"
	"container/list"
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
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
	//数据不完整，直接丢弃
	if len(data) < 10 {
		return ""
	}

	startIndex := 0
	if data[0] == 0x20 {
		startIndex = 1
	}
	for startIndex < len(data)-1 {
		// 检查数据帧起始字节
		if data[startIndex] == 0x00 && data[startIndex+1] == 0x00 {
			if len(data) > startIndex+2 {
				// 根据第三个字节确定数据包长度
				var packetLength int
				switch data[startIndex+2] {
				case 0x05:
					packetLength = 10
				case 0x08:
					packetLength = 13
				default:
					startIndex++
					continue
				}

				// 检查是否有足够的字节来完成这个数据包
				if startIndex+packetLength <= len(data) {
					return handlePacket(data[startIndex : startIndex+packetLength])
					//break
				} else {
					// 数据不完整，跳出循环等待更多数据
					break
				}
			}
		} else {
			startIndex++
		}
	}

	return ""
}
func handlePacket(packet []byte) string {

	var no int64 = 0
	if len(packet) == 10 {
		no = int64(packet[4])<<24 + int64(packet[5])<<16 + int64(packet[6])<<8 + int64(packet[7])
	}
	if len(packet) == 13 {
		no = int64(packet[7])<<24 + int64(packet[8])<<16 + int64(packet[9])<<8 + int64(packet[10])
	}
	msg := common.Substring("0000000000"+strconv.FormatInt(no, 10), len(strconv.FormatInt(no, 10)), 10)
	log.Printf("接收到刷卡信息: %v-%s", packet, msg)
	return msg

}
