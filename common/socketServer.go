package common

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
)

// Packet 定义数据包结构
type Packet struct {
	Conn    net.Conn
	Size    int
	Command string
	Content string
}

func RunServer(packetChan chan<- Packet, address string, port int) {

	// 创建监听器
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		zap.S().Errorln("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	zap.S().Infoln("Server started. Listening on", address, "port", port)

	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			zap.S().Errorln("Error accepting connection:", err.Error())
			return
		}

		// 启动一个 goroutine 处理客户端连接
		go handleClient(conn, packetChan)
	}
}

func handleClient(conn net.Conn, packetChan chan<- Packet) {
	zap.S().Infoln("Client connected:", conn.RemoteAddr().String())
	defer conn.Close()

	// 创建一个缓冲读取器
	reader := bufio.NewReader(conn)

	for {
		// 读取前10个字节，获取数据包的大小
		sizeBytes, err := reader.Peek(10)
		if err != nil {
			if err == io.EOF {
				zap.S().Infoln("Client disconnected:", conn.RemoteAddr().String())
				return
			}
			zap.S().Errorln("Error reading packet size:", err.Error())
			return
		}

		// 解析数据包大小
		packetSize, err := strconv.Atoi(string(sizeBytes))
		if err != nil {
			zap.S().Errorln("Error parsing packet size:", err.Error())
			return
		}

		packetSize = packetSize
		// 读取整个数据包
		packetData := make([]byte, packetSize)
		_, err = io.ReadFull(reader, packetData)
		if err != nil {
			zap.S().Errorln("Error reading packet content:", err.Error())
			return
		}

		// 解析数据包内容和指令
		packetCommand := string(packetData[10:14])
		packetContent := string(packetData[14:])

		// 创建数据包对象
		packet := Packet{
			Conn:    conn,
			Size:    packetSize,
			Command: packetCommand,
			Content: packetContent,
		}

		// 将数据包发送到通道
		packetChan <- packet
	}
}

func SendPacket(conn net.Conn, command string, content string) {
	// 组装数据包内容
	packetContent := command + content

	// 计算数据包大小
	packetSize := len(packetContent) + 10

	// 补充0，使得packetSizeStr的长度达到10个字符
	packetSizeStr := fmt.Sprintf("%010d", packetSize)

	// 拼接数据包
	packetData := packetSizeStr + packetContent

	// 发送数据包
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(packetData)
	if err != nil {
		zap.S().Errorln("Error sending packet:", err.Error())
		return
	}
	writer.Flush()

	fmt.Println("Packet sent:")
	fmt.Println("Size:", packetSize)
	fmt.Println("Command:", command)
	fmt.Println("Content:", content)
}
