package runtime

import (
	"TmdtServer/common"
	"github.com/jacobsa/go-serial/serial"
	"go.uber.org/zap"
	"io"
	"log"
	"time"
)

type UartGpio struct {
	Input_Status1 int //输入状态 -1 不可用 0 断开  1 连接
}

func (ug *UartGpio) Run() {

	pn := common.Config.GetString("ecds.uart1")
	if pn == "" {
		ug.Input_Status1 = -1
		zap.S().Errorln("未初始化物理开关接口！")
		return
	}

	options := serial.OpenOptions{
		PortName:              pn,
		BaudRate:              19200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	p, err := serial.Open(options)
	if err != nil {
		zap.S().Errorln("serial.Open: ", err)
		return
	}

	go ug.Uart_Write(p)

}

func (ug *UartGpio) Uart_Write(p io.ReadWriteCloser) {

	for {
		//向串口写数据标志
		_, err := p.Write([]byte{0x88})
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}

		//读出串口写的数据标志
		buf := make([]byte, 1)
		_, err = p.Read(buf)
		if err != nil {
			ug.Input_Status1 = 0
			if err != io.EOF {
				log.Print("Error reading from serial port: ", err)
			}
		} else {
			if buf[0] == 0x88 {
				ug.Input_Status1 = 1
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
