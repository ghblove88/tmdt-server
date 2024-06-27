package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var deviceIDs = []uint32{1001, 1002, 1003, 1004, 1005, 1006}

func main() {
	serverAddress := "localhost:8866"

	for {
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			fmt.Println("Error connecting, retrying:", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Connected to server")
		sendData(conn)
		conn.Close()
		fmt.Println("Disconnected from server, retrying...")
		time.Sleep(5 * time.Second)
	}
}

func sendData(conn net.Conn) {
	// Create and send data packets every 2 seconds
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		deviceID := deviceIDs[rand.Intn(len(deviceIDs))]
		temp1 := generateTemperature()
		temp2 := generateTemperature()
		temp3 := generateTemperature()

		packet := createDataPacket(deviceID, temp1, temp2, temp3)
		_, err := conn.Write(packet)
		if err != nil {
			fmt.Println("Error sending data:", err.Error())
			return
		}
		// 打印字节数组为16进制格式，每个字节用空格分隔
		for _, b := range packet {
			fmt.Printf("%02x ", b)
		}
		fmt.Printf("\n")
		fmt.Printf("Sent packet: DeviceID=%d, Temp1=%.1f, Temp2=%.1f, Temp3=%.1f\n",
			deviceID, temp1, temp2, temp3)
	}
}

// generateTemperature generates a random temperature value between 36.0 and 37.5
func generateTemperature() float32 {
	return 36.0 + rand.Float32()*(37.5-36.0)
}

// createDataPacket creates a data packet based on the provided parameters
func createDataPacket(deviceID uint32, temp1, temp2, temp3 float32) []byte {
	header1 := byte(0x5A)
	header2 := byte(0xA5)
	command := byte(0x81)

	deviceIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceIDBytes, deviceID)

	temp1Bytes := convertTemperatureToBytes(temp1)
	temp2Bytes := convertTemperatureToBytes(temp2)
	temp3Bytes := convertTemperatureToBytes(temp3)

	data := []byte{
		header1, header2, command,
		deviceIDBytes[0], deviceIDBytes[1], deviceIDBytes[2], deviceIDBytes[3],
		temp1Bytes[0], temp1Bytes[1],
		temp2Bytes[0], temp2Bytes[1],
		temp3Bytes[0], temp3Bytes[1],
	}

	checksum := calculateChecksum(data)
	data = append(data, checksum)
	//data = []byte{0x5a, 0xA5, 0x81, 0x00, 0x00, 0x03, 0xEB,0x01, 0x72, 0x01, 0x74, 0x01, 0x6E, 0xc5}
	return data
}

// convertTemperatureToBytes converts a temperature value to high and low bytes
func convertTemperatureToBytes(temp float32) [2]byte {
	tempValue := int(temp * 10)
	highByte := byte(tempValue >> 8)
	lowByte := byte(tempValue & 0xFF)
	return [2]byte{highByte, lowByte}
}

// calculateChecksum calculates the checksum for a data packet
func calculateChecksum(data []byte) byte {
	sum := 0
	for _, b := range data {
		sum += int(b)
	}
	return byte(sum & 0xFF)
}
