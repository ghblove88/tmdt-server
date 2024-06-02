package runtime

import (
	"TmdtServer/models"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var (
	G_Operator_Info Operator_Info
	G_Doctor_Info   Doctor_Info
	G_Device_Info   Device_Info
	G_Sound_Play    Sound_Play
	G_Socket_Reader Socket_Reader
	G_UartGpio      UartGpio

	G_Operator_Current string // 记录当前操作员
)

type Runtime struct {
}

func (k *Runtime) Run() {

	k.Init()

	//go k.Process()

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

	G_Socket_Reader = Socket_Reader{}
	G_Socket_Reader.Run()

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
