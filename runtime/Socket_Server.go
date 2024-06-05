package runtime

import (
	"TmdtServer/common"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"
)

// DataPacket represents the structure of the received data packet
type DataPacket struct {
	Header1   byte
	Header2   byte
	Command   byte
	DeviceID  [4]byte
	Temp1High byte
	Temp1Low  byte
	Temp2High byte
	Temp2Low  byte
	Temp3High byte
	Temp3Low  byte
	Checksum  byte
}

// DeviceData represents the structure to store device data
type DeviceData struct {
	DeviceID uint32
	Temp1    float32
	Temp2    float32
	Temp3    float32
}

// NewDataPacket creates a new DataPacket from a byte slice
func NewDataPacket(buf []byte) (*DataPacket, error) {
	if len(buf) != 14 {
		return nil, fmt.Errorf("invalid buffer size")
	}
	return &DataPacket{
		Header1:   buf[0],
		Header2:   buf[1],
		Command:   buf[2],
		DeviceID:  [4]byte{buf[3], buf[4], buf[5], buf[6]},
		Temp1High: buf[7],
		Temp1Low:  buf[8],
		Temp2High: buf[9],
		Temp2Low:  buf[10],
		Temp3High: buf[11],
		Temp3Low:  buf[12],
		Checksum:  buf[13],
	}, nil
}

// Validate checks the checksum of the data packet
func (p *DataPacket) Validate() bool {
	sum := p.Header1 + p.Header2 + p.Command +
		p.DeviceID[0] + p.DeviceID[1] + p.DeviceID[2] + p.DeviceID[3] +
		p.Temp1High + p.Temp1Low +
		p.Temp2High + p.Temp2Low +
		p.Temp3High + p.Temp3Low

	return byte(sum&0xFF) == p.Checksum
}

// ToDeviceData converts a DataPacket to a DeviceData structure
func (p *DataPacket) ToDeviceData() DeviceData {
	deviceID := binary.BigEndian.Uint32(p.DeviceID[:])
	temp1 := float32(int(p.Temp1High)<<8|int(p.Temp1Low)) / 10.0
	temp2 := float32(int(p.Temp2High)<<8|int(p.Temp2Low)) / 10.0
	temp3 := float32(int(p.Temp3High)<<8|int(p.Temp3Low)) / 10.0
	return DeviceData{
		DeviceID: deviceID,
		Temp1:    temp1,
		Temp2:    temp2,
		Temp3:    temp3,
	}
}

// ResponsePacket represents the structure of the response data packet
type ResponsePacket struct {
	Header1  byte
	Header2  byte
	Command  byte
	Data     [9]byte
	Checksum byte
}

// NewResponsePacket creates a new ResponsePacket based on the validity of the received data
func NewResponsePacket(valid bool) *ResponsePacket {
	response := &ResponsePacket{
		Header1: 0x5A,
		Header2: 0xA5,
		Command: 0x01,
		Data:    [9]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
	}
	if valid {
		response.Data[0] = 0x66
	} else {
		response.Data[0] = 0x73
	}

	// Calculate checksum
	response.CalculateChecksum()

	return response
}

// CalculateChecksum calculates the checksum for the response packet
func (p *ResponsePacket) CalculateChecksum() {
	sum := p.Header1 + p.Header2 + p.Command
	for _, b := range p.Data {
		sum += b
	}
	p.Checksum = byte(sum & 0xFF)
}

// ToBytes converts the ResponsePacket to a byte slice
func (p *ResponsePacket) ToBytes() []byte {
	buf := make([]byte, 12)
	buf[0] = p.Header1
	buf[1] = p.Header2
	buf[2] = p.Command
	copy(buf[3:12], p.Data[:])
	buf[11] = p.Checksum
	return buf
}

// SocketServer represents the server that handles incoming connections and stores the latest temperature records
type SocketServer struct {
	listener  net.Listener
	dataMutex sync.Mutex
	DataMap   map[uint32]DeviceData
	lastSeen  map[uint32]time.Time
}

func NewServer() (*SocketServer, error) {
	service := common.Config.GetString("socket_server.address") + ":" + common.Config.GetString("socket_server.port")
	listener, err := net.Listen("tcp", service)
	if err != nil {
		return nil, err
	}

	server := &SocketServer{
		listener: listener,
		DataMap:  make(map[uint32]DeviceData),
		lastSeen: make(map[uint32]time.Time),
	}

	go server.cleanupExpiredData()

	return server, nil
}

func (s *SocketServer) Start() {
	zap.S().Infoln("Socket_Reader Running:", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *SocketServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 14)
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			return
		}

		packet, err := NewDataPacket(buf)
		if err != nil {
			fmt.Println("Error parsing data packet:", err.Error())
			continue
		}

		valid := packet.Validate()
		deviceData := packet.ToDeviceData()

		s.updateDataMap(deviceData)

		//time.Sleep(100 * time.Millisecond)
		response := NewResponsePacket(valid)
		_, err = conn.Write(response.ToBytes())
		if err != nil {
			fmt.Println("Error writing response:", err.Error())
			return
		}
	}
}

func (s *SocketServer) updateDataMap(data DeviceData) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	s.DataMap[data.DeviceID] = data
	s.lastSeen[data.DeviceID] = time.Now()
	if *common.TestMode {
		fmt.Printf("Updated data map: %v\n", s.DataMap)
	}
}
func (s *SocketServer) cleanupExpiredData() {
	for {
		time.Sleep(1 * time.Minute)
		s.dataMutex.Lock()
		now := time.Now()
		for deviceID, lastSeen := range s.lastSeen {
			if now.Sub(lastSeen) > 1*time.Minute {
				delete(s.DataMap, deviceID)
				delete(s.lastSeen, deviceID)
				fmt.Printf("Removed device %d due to inactivity\n", deviceID)
			}
		}
		s.dataMutex.Unlock()
	}
}
