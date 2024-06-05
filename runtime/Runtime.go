package runtime

import (
	"TmdtServer/common"
	"TmdtServer/models"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
)

var (
	G_Operator_Info Operator_Info
	G_Doctor_Info   Doctor_Info
	G_Device_Info   Device_Info
	G_Sound_Play    Sound_Play
	G_UartGpio      UartGpio
	G_SocketServer  *SocketServer

	RecordTimer *common.TimerTask // 定时器 := NewTimer()
)

type Runtime struct {
}

func (k *Runtime) Run() {

	k.Init()

	//启动实时运行器
	RecordTimer.Start(10*time.Second, 1*time.Minute, k.RecordTemp)

}

func (k *Runtime) Init() {
	zap.S().Infoln("Database Init......")

	G_Operator_Info = Operator_Info{Operator: make(map[string]STRUCT_OPERATOR_INFO)}
	G_Operator_Info.PullDB()

	G_Doctor_Info = Doctor_Info{Doctor: make(map[string]models.DoctorInfo)}
	G_Doctor_Info.PullDB()

	G_Device_Info = Device_Info{Device: make(map[string]models.DeviceInfo)}
	G_Device_Info.PullDB()

	G_UartGpio = UartGpio{}

	zap.S().Infoln("Database Init Done!")

	G_Sound_Play = Sound_Play{}
	G_Sound_Play.Run()

	G_SocketServer, _ = NewServer()
	G_SocketServer.Start()

	G_UartGpio.Run()
}

// PlayStepVoice 播放指定语音
func (k *Runtime) PlayStepVoice(nVoice_NO string, nStepName_voice []string, Operating string) {

	voice := []string{"wav/ding.wav"}
	for _, v := range nVoice_NO {
		voice = append(voice, "wav/N"+string(v)+".mp3")
	}
	voice = append(voice, "wav/haojing.mp3")
	for _, val := range nStepName_voice {
		voice = append(voice, "wav/"+val+".mp3")
	}

	voice = append(voice, Operating)

	G_Sound_Play.Add_Sound_Queue(voice)
}

func (k *Runtime) RecordTemp() {

}
