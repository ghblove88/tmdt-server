package runtime

import (
	"TmdtServer/common"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"go.uber.org/zap"
	"io"
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
	go sr.SocketStart()
}

func (sr *Socket_Reader) SocketStart() {
	service := common.Config.GetString("socket_server.address") + ":" + common.Config.GetString("socket_server.port")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		zap.S().Errorln("Fatal error: %s", zap.Error(err))
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		zap.S().Errorln("Fatal error: %s", zap.Error(err))
		os.Exit(1)
	}

	zap.S().Infoln("Socket_Reader Running:", tcpAddr)
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
	//conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	var TimeOut int64 = 60 * 10
	Last_Heartbeat := time.Now().Unix()
	defer conn.Close()

	log.Println("From:" + conn.RemoteAddr().String() + " connection is successful。")
	for {
		//连续10分钟没有心跳，就退出
		if time.Now().Unix()-Last_Heartbeat >= TimeOut {
			break
		}
		data, err := readFully(conn)
		if len(data) != 8 {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err != nil {
			break
		}

		if (data[2]&0xff) == 0x00 && (data[3]&0xff) == 0x00 && (data[4]&0xff) == 0x00 && (data[5]&0xff) == 0x00 {
			log.Println(" time: " + string(time.Now().String()) + ", Received Heartbeat: " + conn.RemoteAddr().String())
			Last_Heartbeat = time.Now().Unix()
			continue
		}
		var no int32 = 0
		no_byte := []byte{data[2] & 0xff, data[3] & 0xff, data[4] & 0xff, data[5] & 0xff}
		//no_byte :=[]byte{0x00,0x1e,0x84,0xe4}
		binary.Read(bytes.NewBuffer(no_byte), binary.BigEndian, &no)
		msg := common.Substring("0000000000"+strconv.Itoa(int(no)), len(strconv.Itoa(int(no))), 10)

		remoteip := strings.Split(conn.RemoteAddr().String(), ":")
		G_Socket_Reader.Readid_queue.Push_ReadID_Queue(STRUCT_READID_MSG{remoteip[0], msg})
		log.Println(" time: " + time.Now().String() + ", From:" + conn.RemoteAddr().String() + " Receive credit card: " + msg)

		Last_Heartbeat = time.Now().Unix()
	}
	log.Println("From:" + conn.RemoteAddr().String() + " disconnect。")
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()
	result := bytes.NewBuffer(nil)
	var buf [8]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if n == 8 {
			break
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}
