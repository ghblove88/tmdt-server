package runtime

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	G_TimePlan_Info TimePlan_Info
	G_Program_Info  Program_Info
	G_Operator_Info Operator_Info
	G_Doctor_Info   Doctor_Info
	G_Device_Info   Device_Info
	G_Reader_Info   Reader_Info
	G_Sound_Play    Sound_Play
	G_Socket_Reader Socket_Reader
	G_UartGpio      UartGpio

	G_Operator_Current       string             // 记录当前操作员
	G_Processing_Ant         map[string]AntInfo // 进行中的洗消
	G_Processing_Ant_RWMutex *sync.RWMutex      // 进行中的洗消读写锁
	Sequence_key             []string           // 记录当前进行中的洗消 内窥镜 卡号

	G_Liquid_Expired_Item  []string // 记录洗消液快过期的信息
	G_Liquid_Expired_Count int      // 记录洗消液快过期的数量

	G_Standard_Step     map[string]int // 洗消标准步骤
	G_Standard_Step_Key int            // 洗消标准步骤的关键步骤(这一步骤之后的步骤，不能作为开始步骤)

	G_FreeRule          bool // 是否使用自由规则
	G_AutoEnd           bool // #末洗 指定 分钟到时后自动结束
	G_DryingAutoEnd     bool // #干燥 指定 分钟到时后自动结束
	G_EcmAutoEnd        bool // #洗镜子机 预设时间到时后自动结束
	G_RinseTimeLimit    int  //#漂洗 限定时长 单位 秒
	G_UnwashedTimeLimit int  //#末洗 设置时长 单位 秒
	G_DryingTimeLimit   int  //#干燥 设置时长 单位 秒

	G_ShowLastPatient bool //在洗消监控界面，显示洗消中内镜上一次使用患者名字

	G_RecordAfterEndoscopy bool //使用后记录模式，先记录工作站检查病人信息，再清洗镜子是绑定清洗信息

	G_AddFixedSteps bool // 开始洗消前 自动添加 固定 洗消步骤

)

type Runtime struct {
}

func (k *Runtime) Run() {

	k.Init()

	go k.Process()

}

func (k *Runtime) Init() {
	zap.S().Infoln("Database Init......")
	G_TimePlan_Info = TimePlan_Info{Timeplan: make(map[string]STRUCT_TIMEPLAN_INFO)}
	G_TimePlan_Info.PullDB()

	G_Program_Info = Program_Info{Program: make(map[string]STRUCT_PROGRAM_INFO)}
	G_Program_Info.PullDB()

	G_Operator_Info = Operator_Info{Operator: make(map[string]STRUCT_OPERATOR_INFO)}
	G_Operator_Info.PullDB()

	G_Doctor_Info = Doctor_Info{Doctor: make(map[string]models.DoctorInfo)}
	G_Doctor_Info.PullDB()

	G_Device_Info = Device_Info{Device: make(map[string]models.DeviceInfo)}
	G_Device_Info.PullDB()

	G_UartGpio = UartGpio{}

	G_Reader_Info = Reader_Info{}
	G_Reader_Info.Init()
	zap.S().Infoln("Database Init Done!")

	G_Processing_Ant = make(map[string]AntInfo)
	G_Processing_Ant_RWMutex = new(sync.RWMutex)

	// 初始化自由规则
	G_FreeRule = common.Config.GetBool("ecds.freerule")
	zap.S().Infoln("自由规则：", G_FreeRule)

	G_AutoEnd = common.Config.GetBool("ecds.autoend")
	zap.S().Infoln("末洗 指定 分钟后自动结束：", G_AutoEnd)

	G_DryingAutoEnd = common.Config.GetBool("ecds.dryingautoend")
	zap.S().Infoln("干燥 指定 分钟后自动结束：", G_DryingAutoEnd)

	G_UnwashedTimeLimit = common.Config.GetInt("ecds.unwashedtimelimit")
	zap.S().Infoln("末洗限定时长:", G_UnwashedTimeLimit)

	G_RinseTimeLimit = common.Config.GetInt("ecds.rinsetimelimit")
	zap.S().Infoln("漂洗限定时长:", G_RinseTimeLimit)

	G_DryingTimeLimit = common.Config.GetInt("ecds.dryingtimelimit")
	zap.S().Infoln("干燥限定时长:", G_DryingTimeLimit)

	G_EcmAutoEnd = common.Config.GetBool("ecds.ecmautoend")
	zap.S().Infoln("洗镜子机，预设时间到时后自动结束：", G_EcmAutoEnd)

	G_ShowLastPatient = common.Config.GetBool("ecds.showlastpatient")
	zap.S().Infoln("在洗消监控界面，是否显示上一次使用患者名字：", G_ShowLastPatient)

	G_AddFixedSteps = common.Config.GetBool("ecds.addfixedsteps")
	zap.S().Infoln("预处理步骤，相应操作记录：", G_AddFixedSteps)

	G_RecordAfterEndoscopy = common.Config.GetBool("ecds.recordafterendoscopy")
	if G_RecordAfterEndoscopy { //如果使用后记录模式，就自动关闭 G_ShowLastPatient
		G_ShowLastPatient = false
	}
	zap.S().Infoln("使用后记录模式：", G_RecordAfterEndoscopy)

	// 洗消标准步骤初始化
	G_Standard_Step = make(map[string]int)
	standardStep := common.Config.GetString("ecds.standardstep")
	ss := strings.Split(standardStep, ",")
	for i := 0; i < len(ss); i++ {
		G_Standard_Step[ss[i]] = i + 1
		if ss[i] == "消毒" {
			G_Standard_Step_Key = i + 1
		}
	}
	zap.S().Infoln(G_Standard_Step, G_Standard_Step_Key)

	G_Sound_Play = Sound_Play{}
	G_Sound_Play.Run()

	G_Socket_Reader = Socket_Reader{}
	G_Socket_Reader.Run()

	G_UartGpio.Run()

	//启动定期检测时间线程
	_, _ = NewRegularEvents()

	G_Liquid_Expired_Count = -1

}

func (k *Runtime) Process() {
	for {
		readID, err := G_Socket_Reader.Pull_ReadID_Queue()
		if err != nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		G_Processing_Ant_RWMutex.RLock()
		ant := G_Processing_Ant[readID.MSG]
		G_Processing_Ant_RWMutex.RUnlock()

		if ant.AntID != "" { //卡号 是正在进行洗消的内窥镜
			if res := G_Device_Info.GetAt(readID.MSG); res.EndoscopeNumber != "" {

				res := k.execute(res, readID, &ant)
				//只有不是结束头 才保存ant 状态
				if res != 1 {
					G_Processing_Ant_RWMutex.Lock()
					G_Processing_Ant[readID.MSG] = ant
					G_Processing_Ant_RWMutex.Unlock()
				}
			}
			continue

		} else if res := G_Operator_Info.GetAt(readID.MSG); res.Name != "" { //当前为操作员卡

			G_Operator_Current = res.Name
			G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czyyqr.mp3"})
			log.Println("时间：" + time.Now().String() + " 操作员：" + G_Operator_Current + " 确认为当前操作员。")
			continue

		} else if res := G_Device_Info.GetAt(readID.MSG); res.EndoscopeNumber != "" { //为内窥镜卡 新刷

			reader_res := G_Reader_Info.Find(readID.IP)

			if reader_res.Name == "结束" {
				G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
				log.Println("时间：" + time.Now().String() + " 操作错误，没有进行洗消的内窥镜刷了结束！")
				continue
			}

			if G_Operator_Current == "" && reader_res.Name != "送达" && reader_res.Name != "预处理" {
				G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/qsczyk.mp3"})
				log.Println("时间：" + time.Now().String() + " 没有指定操作员！")
				continue
			}

			if reader_res.Name == "送达" {
				sql_str := "insert into endoscope_sign(endoscope_number, time, ant_number) values (" +
					"'" + readID.MSG + "'," +
					"'" + time.Now().Format("2006-01-02 15:04:05") + "','')"
				engine, _ := models.NewOrm()
				db := engine.Exec(sql_str)
				if db.Error != nil {
					log.Println(db.Error)
					continue
				}
				G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav"})
				log.Println("时间：" + time.Now().String() + readID.MSG + " 内窥镜送入洗消中心！")
				continue
			}

			if reader_res.Name == "存镜" {
				AccessEndoscope(readID.MSG, G_Operator_Current, "存入")
				G_Operator_Current = ""
				continue
			}

			if reader_res.Name == "取镜" {
				AccessEndoscope(readID.MSG, G_Operator_Current, "取出")
				G_Operator_Current = ""
				continue
			}

			if reader_res.Name == "预处理" {
				PreprocessEndoscope(readID.MSG, G_Operator_Current, readID.IP, "预处理")
				G_Operator_Current = ""
				continue
			}

			k.executeNew(res, readID) //创建新洗消进程
			G_Operator_Current = ""
			continue
		}
	}
}

func (k *Runtime) executeNew(device models.DeviceInfo, msg STRUCT_READID_MSG) {

	res := G_Reader_Info.Find(msg.IP)
	if res.Name == "" {
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
		return
	}

	// >0 代表这个步骤是在 标准步骤中 >G_Standard_Step_Key 代表这个步骤超过了消毒 （消毒代码为 G_Standard_Step_Key)
	if G_Standard_Step[res.Name] > 0 && G_Standard_Step[res.Name] > G_Standard_Step_Key {
		log.Print("错误:在限定步骤顺序的情况下，从消毒之后的步骤开始！")
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
		return
	}

	// 在标准步骤限制下 res.Name == "漂洗" 不允许从 消毒之前的 一个步骤开始（"漂洗"）
	if len(G_Standard_Step) > 0 && res.Name == "漂洗" {
		log.Print("错误:在限定步骤顺序的情况下，从 漂洗 步骤开始！")
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
		return
	}

	//向数据库里写入记录
	ant := AntInfo{}
	ant_number_res := ant.Init(device.EndoscopeNumber, G_Operator_Current, device.EndoscopeType, device.EndoscopeInfo)
	if !ant_number_res {
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/sjcw.mp3"})
		return
	}

	// 床旁预处理
	if G_AddFixedSteps == true && (res.Name == "测漏" || res.Name == "清洗") {

		engine, _ := models.NewOrm()
		var pl models.PreprocessLog
		engine.Where("endoscope_number=? and number = ''", device.EndoscopeNumber).Order("id desc").First(&pl)

		if pl.Id > 0 {
			k.AddFixedSteps(&ant, "预处理", (int)(time.Now().Sub(pl.Time).Seconds()), pl.Ip)

			db := engine.Model(&models.PreprocessLog{}).Where("id=?", pl.Id).Update("number", ant.AntID)
			if db.Error != nil {
				log.Println(db.Error)
			}
		}
	}

	//如果新开始的步骤是 测漏步骤 就增加一条 测漏记录
	if res.Name == "测漏" {
		ant.InitLeakDetectionRecord(ant.AntID)
	}

	var time_long int
	if res.Name == "清洗" {
		time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 0)
	} else if res.Name == "消毒" {
		if ant.Diseases == 2 { // 有传染病 增加消毒时间
			time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 2)
		} else {
			time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 1)
		}
	} else if res.Name == "漂洗" {
		//如果配置文件中，漂洗限时，则限制漂洗结束时间
		if G_RinseTimeLimit > 0 {
			time_long = G_RinseTimeLimit
		}
	} else if res.Name == "干燥" {
		//如果配置文件中，干燥限时，则限制干燥结束时间
		if G_DryingTimeLimit > 0 {
			time_long = G_DryingTimeLimit
		}
	}
	if res.Name == "末洗" {
		time_long = G_UnwashedTimeLimit
	} else if strings.Contains(res.Name, "洗镜机") {
		//每日第一次使用洗净机时，调用浸泡程序
		if k.IsFirstTimeOfDay(msg) {
			reg := regexp.MustCompile(`[^\p{Han}]+`)
			res.Name = "洗镜机浸泡" + reg.FindString(res.Name)
		}
		time_long = G_Program_Info.Program[res.Name].TotalCostTime
	}
	ant.CurStepName = res.Name
	ant.CurStepIP = msg.IP
	ant.CurStepTimeLong = time_long
	ant.CurStepStartTime = time.Now()

	// 是否显示上一次病人姓名 2022-10-04
	if G_ShowLastPatient {
		go GetLastPatient(msg.MSG)
	}

	// 后记录模式，开始洗消同时，获取病人信息，实现绑定 2022-10-12
	if G_RecordAfterEndoscopy {
		go RecordAfterEndoscopy(msg.MSG)
	}

	G_Processing_Ant_RWMutex.Lock()
	G_Processing_Ant[msg.MSG] = ant
	G_Processing_Ant_RWMutex.Unlock()
	//使用list 记录当前进行中的洗消记录
	Sequence_key = append(Sequence_key, msg.MSG)

	k.PlayStepVoice(ant.VoiceNo, res.Voice, "wav/kaishi.mp3")
	if time_long > 0 {
		go CountdownPlayStepTimeout(time_long, res.Voice, ant.VoiceNo)
	}

	// 末洗 指定 分钟到时后自动结束
	if G_AutoEnd && res.Name == "末洗" {
		go CountdownAutoEnd(time_long, ant.No)
	}

	// 干燥 指定 分钟到时后自动结束
	if G_DryingAutoEnd && res.Name == "干燥" {
		go CountdownAutoEnd(time_long, ant.No)
	}

	// 洗镜子机 预设时间到时后自动结束
	if G_EcmAutoEnd && strings.Contains(res.Name, "洗镜机") {
		go CountdownAutoEnd(time_long, ant.No)
	}

	// 更新这次洗消 内窥镜 的 送达记录表
	sqlStr := "update endoscope_sign set ant_number='" + ant.AntID + "' where endoscope_number='" + ant.No +
		"' and ant_number='' order by time desc limit 1"
	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}
	db := engine.Exec(sqlStr)
	if db.Error != nil {
		log.Println(db.Error)
	}
}

// GetLastPatient 是否显示上一次病人姓名 2022-10-04
func GetLastPatient(no string) {

	engine, err := models.NewOrm()

	if err != nil {
		println(err.Error())
		return
	}

	ant := new(models.Ant)
	//_, err = engine.Where("patient_name <> ''").And("endoscope_number = ?", no).
	//	And("begin_time <> end_time").Desc("ID").Limit(1, 0).Get(ant)
	engine.Where("endoscope_number = ? and begin_time <> end_time", no).Order("id desc").First(ant)

	if len(ant.PatientName) <= 0 {
		return
	}

	G_Processing_Ant_RWMutex.RLock()
	temp := G_Processing_Ant[no]
	G_Processing_Ant_RWMutex.RUnlock()

	temp.PatientName = "^" + ant.PatientName

	G_Processing_Ant_RWMutex.Lock()
	G_Processing_Ant[no] = temp
	G_Processing_Ant_RWMutex.Unlock()
}

// RecordAfterEndoscopy 后记录模式，开始洗消同时，获取病人信息，实现绑定 2022-10-12
func RecordAfterEndoscopy(no string) {

	engine, err := models.NewOrm()

	if err != nil {
		println(err.Error())
		return
	}

	ewb := new(models.EndoscopyWriteback)
	engine.Where("rfidNo = ?", no).Order("id desc").First(ewb)

	if len(ewb.Namepatient) <= 0 || len(ewb.Number) > 0 {
		return
	}

	ewb.Number = G_Processing_Ant[no].AntID
	ewb.BeginTime = G_Processing_Ant[no].CurStepStartTime
	ewb.Bindmirrostate = "1"
	db := engine.Save(ewb)
	if db.Error != nil {
		println(db.Error)
		return
	}

	SqlString := "update ant set patient_name='" + ewb.Namepatient + "' where number='" + G_Processing_Ant[no].AntID + "'"
	engine.Exec(SqlString)

	G_Processing_Ant_RWMutex.RLock()
	temp := G_Processing_Ant[no]
	G_Processing_Ant_RWMutex.RUnlock()

	temp.PatientName = ewb.Namepatient

	G_Processing_Ant_RWMutex.Lock()
	G_Processing_Ant[no] = temp
	G_Processing_Ant_RWMutex.Unlock()
}

func (k *Runtime) execute(device models.DeviceInfo, msg STRUCT_READID_MSG, ant *AntInfo) (res_status int64) {

	res_status = 0

	res := G_Reader_Info.Find(msg.IP)
	if res.Name == "" {
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
		return res_status
	}

	if res.Name == ant.CurStepName { // 刷卡为当前操作步骤

		//执行消毒的过程中，又刷了消毒，即进入灭菌模式，同时标注为传染病
		if res.Name == "消毒" {
			ant.Diseases = 2
			ant.CurStepTimeLong = G_TimePlan_Info.GetAt(device.EndoscopeType, 2)
			AntDiseasesUpdate(ant.AntID, 2)
			k.PlayStepVoice(ant.VoiceNo, []string{}, "wav/mj.mp3")
			return res_status
		}

		k.PlayStepVoice(ant.VoiceNo, []string{}, "wav/zzzxdqbz.mp3")
		log.Println(" time: " + time.Now().Format("2006-01-02 15:04:05") + ": 继续处理已启动的洗消流程，刷卡操作有误，正在执行当前步骤！")
		return res_status
	}

	// 如果当前步骤 + 1 不等于 即将开始的步骤，就是操作错误
	// 步骤列表不为空 && 即将开始的不是洗净机 && 当前步骤是标准中的步骤 && 下一步骤是紧挨着当前步骤的
	if len(G_Standard_Step) > 0 && !strings.Contains(res.Name, "洗镜机") &&
		G_Standard_Step[ant.CurStepName] > 0 && G_Standard_Step[ant.CurStepName]+1 != G_Standard_Step[res.Name] {

		log.Println(time.Now().Format("2006-01-02 15:04:05") + " 错误:在限定步骤顺序的情况下，当前步骤 + 1 不等于 即将开始的步骤！")
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3"})
		return
	}

	// 当前步骤操作时间不足，不能进行下一步操作
	if ant.CurStepTimeLong > 0 && time.Now().Unix()-ant.CurStepStartTime.Unix() < int64(ant.CurStepTimeLong) {
		// && !strings.Contains(msg.IP, "8.8.8")  这个判断是为 ECM洗镜机专用,当ECM洗镜上上报结束时,会用8.8.8.xx做IP
		// 检测到这种IP 时,即使时间不足,也可以结束当前洗镜机操作
		if !strings.Contains(msg.IP, "8.8.8") {
			if ant.CurStepName == "清洗" {
				k.PlayStepVoice(ant.VoiceNo, []string{"qingxi"}, "wav/sjbz.mp3")
			} else if ant.CurStepName == "消毒" {
				k.PlayStepVoice(ant.VoiceNo, []string{"xiaodu"}, "wav/sjbz.mp3")
			} else if ant.CurStepName == "末洗" {
				k.PlayStepVoice(ant.VoiceNo, []string{"moxi"}, "wav/sjbz.mp3")
			} else if ant.CurStepName == "漂洗" {
				k.PlayStepVoice(ant.VoiceNo, []string{"piaoxi"}, "wav/sjbz.mp3")
			} else if strings.Contains(ant.CurStepName, "洗镜机") {
				k.PlayStepVoice(ant.VoiceNo, []string{"xjj"}, "wav/sjbz.mp3")
			}
			log.Println(" time: " + time.Now().Format("2006-01-02 15:04:05") + ":  继续处理已启动的洗消流程，刷卡操作有误，操作时间不足！")
			return res_status
		}
	}
	// 当刷了【结束】头时判断是否符合结束规则
	if res.Name == "结束" {
		if !k.CheckEndRules(ant) {
			G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/qsbz.mp3"})
			log.Println(" time: " + time.Now().Format("2006-01-02 15:04:05") + ":  当前洗消不符合结束规则！")
			return
		}
	}

	// 将当前运行步骤保存到数据库
	k.AddCompleteStep(ant)

	// 当前头是结束
	if res.Name == "结束" {

		// 洗消结束后硬加两个步骤 2023-4-17 ！！！！
		/*
			k.AddCompleteStep(&Ant_Info{ant.Voice_NO, ant.AntID, ant.No, ant.Type, ant.Info,
				ant.Operator, ant.PatientName, ant.DoctorName, ant.Diseases,
				"酒精", "192.168.10.70", time.Now().Add(-time.Second * 4), 4,
				ant.StepList, ant.Step_TimeLongList})

			k.AddCompleteStep(&Ant_Info{ant.Voice_NO, ant.AntID, ant.No, ant.Type, ant.Info,
				ant.Operator, ant.PatientName, ant.DoctorName, ant.Diseases,
				"干燥", "192.168.10.80", time.Now().Add(-time.Second * 30), 30,
				ant.StepList, ant.Step_TimeLongList})
		*/

		// 保存数据 结束洗消
		ant.Destroy()

		k.PlayStepVoice(ant.VoiceNo, []string{}, "wav/jieshu.mp3")
		G_Processing_Ant_RWMutex.Lock()
		delete(G_Processing_Ant, msg.MSG)
		G_Processing_Ant_RWMutex.Unlock()

		for i := 0; i < len(Sequence_key); i++ {
			if Sequence_key[i] == msg.MSG {
				Sequence_key = common.Slice_Remove(Sequence_key, i, i+1)
				break
			}
		}
		res_status = 1
		return res_status
	}

	// 开始下一步骤

	// 如果新开始的步骤是 测漏步骤 就增加一条 测漏记录
	if res.Name == "测漏" {
		ant.InitLeakDetectionRecord(ant.AntID)
	}

	var time_long int
	if res.Name == "清洗" {
		time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 0)
	} else if res.Name == "消毒" {
		if ant.Diseases == 2 { // 有传染病 增加消毒时间
			time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 2)
		} else {
			time_long = G_TimePlan_Info.GetAt(device.EndoscopeType, 1)
		}
	} else if res.Name == "末洗" {
		time_long = G_UnwashedTimeLimit
	} else if res.Name == "漂洗" {
		//如果配置文件中，漂洗限时，则限制漂洗结束时间
		if G_RinseTimeLimit > 0 {
			time_long = G_RinseTimeLimit
		}
	} else if res.Name == "干燥" {
		//如果配置文件中，干燥限时，则限制干燥结束时间
		if G_DryingTimeLimit > 0 {
			time_long = G_DryingTimeLimit
		}
	} else if strings.Contains(res.Name, "洗镜机") {
		// 每日第一次使用洗净机时，调用浸泡程序
		if k.IsFirstTimeOfDay(msg) {
			reg := regexp.MustCompile(`[^\p{Han}]+`)
			res.Name = "洗镜机浸泡" + reg.FindString(res.Name)
		}
		time_long = G_Program_Info.Program[res.Name].TotalCostTime
	}
	ant.CurStepName = res.Name
	ant.CurStepIP = msg.IP
	ant.CurStepTimeLong = time_long
	ant.CurStepStartTime = time.Now()

	k.PlayStepVoice(ant.VoiceNo, res.Voice, "wav/kaishi.mp3")
	if time_long > 0 {
		go CountdownPlayStepTimeout(time_long, res.Voice, ant.VoiceNo)
	}
	// 末洗 指定 分钟到时后自动结束
	if G_AutoEnd && res.Name == "末洗" {
		go CountdownAutoEnd(time_long, ant.No)
	}

	// 干燥 指定 分钟到时后自动结束
	if G_DryingAutoEnd && res.Name == "干燥" {
		go CountdownAutoEnd(time_long, ant.No)
	}

	// 洗镜子机 预设时间到时后自动结束
	if G_EcmAutoEnd && strings.Contains(res.Name, "洗镜机") {
		go CountdownAutoEnd(time_long, ant.No)
	}

	return res_status
}

func (k *Runtime) UpdateUsageCount(ip string) {

	SqlString := "select id from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=(select id from liquid_device where device_info = '" + ip + "') ORDER BY id desc limit 0,1"
	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}
	var rowID int
	engine.Raw(SqlString).Scan(&rowID)
	if rowID > 0 {
		db := engine.Model(&models.LiquidDetectionRecord{}).Where("id = ?", rowID).Update("usage_count", gorm.Expr("usage_count - ?", 1))
		if db.Error != nil {
			log.Println("UpdateUsageCount:", db.Error)
		}
	}
	log.Println("UpdateUsageCount:", ip)
}

func (k *Runtime) AddCompleteStep(ant *AntInfo) {

	curTimeLong := time.Now().Unix() - ant.CurStepStartTime.Unix()
	ant.StepList = append(ant.StepList, ant.CurStepName)
	ant.StepTimeLongList = append(ant.StepTimeLongList, int(curTimeLong))

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}

	if strings.Contains(ant.CurStepName, "洗镜机") {
		sri := G_Program_Info.Program[ant.CurStepName]
		var SqlString = []string{}
		for i := 0; i < len(sri.StepList); i++ {
			SqlString = append(SqlString, "insert into ant_step(number,step,stepip, cost_time, washing_machine) values ("+
				"'"+ant.AntID+"',"+
				"'"+sri.StepList[strconv.Itoa(i+1)].Step+"',"+
				"'"+ant.CurStepIP+"',"+
				"'"+strconv.Itoa(sri.StepList[strconv.Itoa(i+1)].Cost_time)+"',"+
				"'"+ant.CurStepName+"')")

			//消毒步骤使用，就减少一次对应设备消毒液的 可用次数
			if sri.StepList[strconv.Itoa(i+1)].Step == "消毒" {
				k.UpdateUsageCount(ant.CurStepIP)
			}
			//如果有 测漏步骤 就增加一条 测漏记录
			if ant.CurStepName == "测漏" {
				ant.InitLeakDetectionRecord(ant.AntID)
			}
		}
		for _, val := range SqlString {
			db := engine.Exec(val)
			if db.Error != nil {
				log.Println(db.Error)
			}
		}

	} else {
		SqlString := "insert into ant_step(number,step,stepip, cost_time, washing_machine) values (" +
			"'" + ant.AntID + "'," +
			"'" + ant.CurStepName + "'," +
			"'" + ant.CurStepIP + "'," +
			"'" + strconv.FormatInt(curTimeLong, 10) + "','')"

		db := engine.Exec(SqlString)
		if db.Error != nil {
			log.Println(db.Error)
		}
		//消毒步骤使用，就减少一次对应设备消毒液的 可用次数
		if ant.CurStepName == "消毒" {
			k.UpdateUsageCount(ant.CurStepIP)
		}
	}
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

// CountdownPlayStepTimeout 限制操作时长的步骤 开始倒计时 播放时间已到语音
func CountdownPlayStepTimeout(nTimeLong int, stepVoice []string, nVoice_NO string) {

	voice := []string{"wav/ding.wav"}
	for _, v := range nVoice_NO {
		voice = append(voice, "wav/N"+string(v)+".mp3")
	}
	voice = append(voice, "wav/haojing.mp3")
	for _, val := range stepVoice {
		voice = append(voice, "wav/"+val+".mp3")
	}
	voice = append(voice, "wav/sjyd.mp3")
	time.Sleep(time.Duration(nTimeLong) * time.Second)

	G_Sound_Play.Add_Sound_Queue(voice)
}

// CountdownAutoEnd 末洗 2分钟到时后自动结束
func CountdownAutoEnd(nTimeLong int, number string) {

	res := G_Reader_Info.FindByName("结束")
	time.Sleep(time.Duration(nTimeLong) * time.Second)
	G_Socket_Reader.Readid_queue.Push_ReadID_Queue(
		STRUCT_READID_MSG{"8.8.8." + res.Type, number})
}

// IsFirstTimeOfDay 判断当前内镜清洗是否为本日第一次
func (k *Runtime) IsFirstTimeOfDay(msg STRUCT_READID_MSG) bool {

	//如果配置文件中 设置 第一次判断为 false 不使用此功能
	isFday := common.Config.GetString("ecds.isfirstsoakofday")
	if isFday == "false" {
		return false
	}

	//物理开关可用
	if G_UartGpio.Input_Status1 != -1 {
		if G_UartGpio.Input_Status1 == 0 {
			return false
		} else if G_UartGpio.Input_Status1 == 1 {
			return true
		}
	}

	Sql_String := "select ID from ant where begin_time <> end_time and endoscope_number='" + msg.MSG +
		"' and begin_time between '" + time.Now().Format("2006-01-02 00:00:00") + "' and '" + time.Now().Format("2006-01-02 15:04:05") + "' ORDER BY ID DESC  limit 0,1 "

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return false
	}
	var antID int
	engine.Raw(Sql_String).Scan(&antID)
	if antID <= 0 {
		return true
	} else {
		return false
	}
}

// CheckEndRules 判断是否符合结束规则
func (k *Runtime) CheckEndRules(ant *AntInfo) bool {

	// 如果是自由规则模式 直接返回真
	if G_FreeRule {
		return true
	}

	// 如果是洗净机，直接返回真
	if strings.Contains(ant.CurStepName, "洗镜机") {
		return true
	}
	// 如果标准列表不为空 就直接返回真
	if len(G_Standard_Step) > 0 {
		return true
	}

	step_value := append(ant.StepList, ant.CurStepName)
	// 规则一
	if len(step_value) == 2 && step_value[0] == "消毒" && step_value[1] == "末洗" {
		return true
	}
	// 规则二
	if len(step_value) >= 4 {
		q1, q2, q3, q4 := false, false, false, false
		for _, str := range step_value {
			if str == "清洗" {
				q1 = true
			}
			if str == "漂洗" {
				q2 = true
			}
			if str == "消毒" {
				q3 = true
			}
			if str == "末洗" {
				q4 = true
			}
			// 如果有步骤是洗净机，直接返回真
			if strings.Contains(str, "洗镜机") {
				return true
			}
		}
		if q1 && q2 && q3 && q4 {
			return true
		} else {
			return false
		}
	}
	return false
}

// 洗消的记录开始前，增加指定步骤，同时更改数据库的开始时间，更新洗消界面显示内容
func (k *Runtime) AddFixedSteps(ant *AntInfo, stepName string, timeLong int, stepIP string) {

	////随机生成100以内的正整数
	//rand.NewSource(time.Now().UnixNano())
	//fmt.Println(rand.Intn(420))
	curTimeLong := timeLong
	ant.StepList = append(ant.StepList, stepName)
	ant.StepTimeLongList = append(ant.StepTimeLongList, timeLong)

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}

	sqlstr := "insert into ant_step(number,step,stepip, cost_time, washing_machine) values (" +
		"'" + ant.AntID + "'," +
		"'" + stepName + "'," +
		"'" + stepIP + "'," +
		"'" + strconv.Itoa(curTimeLong) + "','')"

	db := engine.Exec(sqlstr)
	if db.Error != nil {
		log.Println(db.Error)
	}
	//消毒步骤使用，就减少一次对应设备消毒液的 可用次数
	if ant.CurStepName == "消毒" {
		k.UpdateUsageCount(ant.CurStepIP)
	}

	sqlstr = "update ant set begin_time=DATE_SUB(begin_time,INTERVAL " + strconv.Itoa(curTimeLong) +
		" second) ,end_time = begin_time where number='" + ant.AntID + "'"
	db = engine.Exec(sqlstr)
	if db.Error != nil {
		log.Println(db.Error)
	}
}

// AntDataUpdate 更新洗消记录 病人姓名 医生姓名 是否有传染病
func AntDataUpdate(ant_Number string, patient_name string, doc_name string, diseases string) {

	SqlString := "update ant set endoscope_number = endoscope_number"
	if !strings.Contains(patient_name, "^") {
		SqlString += ",patient_name='" + patient_name + "'"
	}
	if len(doc_name) > 0 {
		SqlString += ",doc_name='" + doc_name + "'"
	}
	if len(diseases) > 0 {
		SqlString += ",diseases='" + diseases + "'"
	}
	SqlString += " where number='" + ant_Number + "'"

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}
	db := engine.Exec(SqlString)
	if db.Error != nil {
		log.Println(db.Error)
	}
}

// AntDiseasesUpdate 更新是否有传染病
func AntDiseasesUpdate(ant_Number string, diseases int) {

	SqlString := fmt.Sprintf("update ant set diseases='%d' where number='%s'", diseases, ant_Number)
	log.Println(SqlString)
	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}
	db := engine.Exec(SqlString)
	if db.Error != nil {
		log.Println(db.Error)
	}
}

func AntDataLeakPart(ant_Number string, leakparts string) {
	temp := ""
	if len(leakparts) <= 0 {
		temp = "conclusion='完好',leakparts=''"
	} else {
		temp = "conclusion='漏水',leakparts='" + leakparts + "'"
	}

	SqlString := "update leak_detection_record set " + temp + " where number='" + ant_Number + "'"
	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return
	}

	db := engine.Exec(SqlString)
	if db.Error != nil {
		log.Println(db.Error)
	}
}

// GetLeakPart 获取指定洗消编号的 测漏数据
func GetLeakPart(ant_Number string) (res string) {
	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}
	var leakDetectionRecord models.LeakDetectionRecord
	engine.Where("number = ?", ant_Number).First(&leakDetectionRecord)
	res = ""
	if leakDetectionRecord.Id > 0 {
		res = leakDetectionRecord.Leakparts
	} else {
		res = "无记录"
	}
	return res
}

// GetAntRecord 获取洗消记录
func GetAntRecord(SqlString string) (res map[string]AntDb_1) {

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}
	var ants []map[string]string
	engine.Raw(SqlString).Scan(&ants)
	res = make(map[string]AntDb_1)

	for i, val := range ants {
		steps := make(map[string]AntStepDb)
		// 2021-10-7 返回数据增加 endoscope_model 通过截取 endoscope_info实现
		endoscope_model := strings.Split(val["endoscope_info"], "/")[0]

		_db := AntDb_1{val["id"], val["number"], val["endoscope_number"],
			val["endoscope_type"], val["operator"], val["patient_name"],
			val["doc_name"], val["diseases"],
			val["begin_time"],
			val["end_time"], val["total_cost_time"],
			val["endoscope_info"], endoscope_model, steps}

		var ant_steps []map[string]string
		engine.Raw("select * from ant_step where number='" + val["number"] + "' ORDER BY id").Scan(&ant_steps)

		for x, val_step := range ant_steps {
			_step := AntStepDb{val_step["id"], "", "", "", "", ""}
			if val_step["number"] != "" {
				_step.S_number = val_step["number"]
			}
			if val_step["step"] != "" {
				_step.S_step = val_step["step"]
				if val["diseases"] == "2" && val_step["step"] == "消毒" {
					_step.S_step = "消毒(灭菌)"
				}
			}
			if val_step["stepip"] != "" {
				_step.S_stepip = val_step["stepip"]
			}
			if val_step["cost_time"] != "" {
				_step.S_cost_time = val_step["cost_time"]
			}
			if val_step["washing_machine"] != "" {
				_step.S_washing_machine = val_step["washing_machine"]
			}
			steps[strconv.Itoa(x+1)] = _step
		}
		res[strconv.Itoa(i+1)] = _db
	}
	return res
}

// GetAntRecordSlic 获取洗消记录 数组
func GetAntRecordSlic(SqlString string) (res []AntDbSlic_1) {

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}

	var ants []map[string]string
	engine.Raw(SqlString).Scan(&ants)
	res = []AntDbSlic_1{}

	for _, val := range ants {
		steps := []AntStepDb{}
		var ant_steps []map[string]string
		engine.Raw("select * from ant_step where number='" + val["number"] + "' ORDER BY id").Scan(&ant_steps)

		for _, val_step := range ant_steps {
			_step := AntStepDb{val_step["id"], "", "", "", "", ""}
			if val_step["number"] != "" {
				_step.S_number = val_step["number"]
			}
			if val_step["step"] != "" {
				_step.S_step = val_step["step"]
				if val["diseases"] == "2" && val_step["step"] == "消毒" {
					_step.S_step = "消毒(灭菌)"
				}
			}
			if val_step["stepip"] != "" {
				_step.S_stepip = val_step["stepip"]
			}
			if val_step["cost_time"] != "" {
				_step.S_cost_time = val_step["cost_time"]
			}
			if val_step["washing_machine"] != "" {
				_step.S_washing_machine = val_step["washing_machine"]
			}
			steps = append(steps, _step)
		}
		// 2021-10-7 返回数据增加 endoscope_model 通过截取 endoscope_info实现
		endoscope_model := strings.Split(val["endoscope_info"], "/")[0]

		_db := AntDbSlic_1{val["id"], val["number"], val["endoscope_number"],
			val["endoscope_type"], val["operator"], val["patient_name"],
			val["doc_name"], val["diseases"],
			val["begin_time"],
			val["end_time"], val["total_cost_time"],
			val["endoscope_info"], endoscope_model, steps}
		res = append(res, _db)
	}
	return res
}

// GetAntRecordSlicUsed 获取洗消记录 数组
func GetAntRecordSlicUsed(SqlString string) (res []AntDbSlic_2) {

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}

	var ants []map[string]string
	engine.Raw(SqlString).Scan(&ants)
	res = []AntDbSlic_2{}

	for _, val := range ants {
		steps := []AntStepDb{}
		var ant_steps []map[string]string
		engine.Raw("select * from ant_step where number='" + val["number"] + "' ORDER BY id").Scan(&ant_steps)

		for _, val_step := range ant_steps {
			_step := AntStepDb{val_step["id"], "", "", "", "", ""}
			if val_step["number"] != "" {
				_step.S_number = val_step["number"]
			}
			if val_step["step"] != "" {
				_step.S_step = val_step["step"]
				if val["diseases"] == "2" && val_step["step"] == "消毒" {
					_step.S_step = "消毒(灭菌)"
				}
			}
			if val_step["stepip"] != "" {
				_step.S_stepip = val_step["stepip"]
			}
			if val_step["cost_time"] != "" {
				_step.S_cost_time = val_step["cost_time"]
			}
			if val_step["washing_machine"] != "" {
				_step.S_washing_machine = val_step["washing_machine"]
			}
			steps = append(steps, _step)
		}
		// 2021-10-7 返回数据增加 endoscope_model 通过截取 endoscope_info实现
		endoscope_model := strings.Split(val["endoscope_info"], "/")[0]
		patientVisitId := ""
		checkReportId := ""

		if len(val["patient_name"]) > 0 {
			var patient_info []map[string]string
			engine.Raw("select * from dongguan_writeback where number='" + val["number"] + "' ORDER BY id desc limit 0,1").Scan(&patient_info)

			if len(patient_info) > 0 {
				patientVisitId = patient_info[0]["patientVisitId"]
				checkReportId = patient_info[0]["checkReportId"]
			}
		}

		_db := AntDbSlic_2{val["id"], val["number"], val["endoscope_number"],
			val["endoscope_type"], val["operator"], val["patient_name"],
			val["doc_name"], val["diseases"],
			val["begin_time"],
			val["end_time"], val["total_cost_time"],
			val["endoscope_info"], endoscope_model,
			patientVisitId, checkReportId, steps}
		res = append(res, _db)
	}
	return res
}

// AccessEndoscope 存取内窥镜操作
func AccessEndoscope(enumber string, operator string, access string) (res bool) {

	sql_str := "select id,endoscope_number,access from store_log where endoscope_number='" + enumber + "' Order By id Desc limit 0 ,1"

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}

	var slog models.StoreLog
	engine.Select("id", "endoscope", "number", "access").Where("endoscope_number = ?", enumber).Order("id desc").First(&slog)

	//计算符合条件 总条数
	//count, _ := o.Raw(sql_str).Values(&ant)

	if access == "存入" {
		if slog.Id > 0 && slog.Access == "存入" {
			G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3", "wav/yicunru.mp3"})
			log.Println("时间：" + time.Now().String() + enumber + " 内窥镜【" + access + "】镜库失败，已经在库！")
			return false
		}
		sql_str = "insert into store_log(endoscope_number, operator,time, access) values (" +
			"'" + enumber + "'," +
			"'" + operator + "'," +
			"'" + time.Now().Format("2006-01-02 15:04:05") + "','存入')"

		db := engine.Exec(sql_str)
		if db.Error != nil {
			log.Println(db.Error)
			return false
		}

		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/cunjing.mp3"})
		log.Println("时间：" + time.Now().String() + enumber + " 内窥镜【" + access + "】镜库！")

	} else if access == "取出" {
		if slog.Id > 0 && slog.Access == "取出" {
			G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/czcw.mp3", "wav/yiquchu.mp3"})
			log.Println("时间：" + time.Now().String() + enumber + " 内窥镜【" + access + "】镜库失败，已经出库！")
			return false
		}
		sql_str = "insert into store_log(endoscope_number, operator,time, access) values (" +
			"'" + enumber + "'," +
			"'" + operator + "'," +
			"'" + time.Now().Format("2006-01-02 15:04:05") + "','取出')"

		db := engine.Exec(sql_str)
		if db.Error != nil {
			log.Println(db.Error)
			return false
		}

		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/qujing.mp3"})
		log.Println("时间：" + time.Now().String() + enumber + " 内窥镜【" + access + "】镜库！")
	}

	return true
}

// PreprocessEndoscope 预处理操做
func PreprocessEndoscope(enumber string, operator string, ip string, access string) (res bool) {

	engine, dberr := models.NewOrm()
	if dberr != nil {
		println(dberr.Error())
		return res
	}

	if access == "预处理" {

		sql_str := "insert into preprocess_log(endoscope_number, operator,time,ip, number) values (" +
			"'" + enumber + "'," +
			"'" + operator + "'," +
			"'" + time.Now().Format("2006-01-02 15:04:05") + "'," +
			"'" + ip + "','')"

		db := engine.Exec(sql_str)
		if db.Error != nil {
			log.Println(db.Error)
			return false
		}

		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav"})
		log.Println("时间：" + time.Now().String() + enumber + " 内窥镜【" + access + "】记录！")

	}

	return true
}
