package V1

import (
	"EcdsServer/common"
	"EcdsServer/runtime"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Ant_Info struct {
	// 语音播报编号
	Voice_NO string `json:"voice"`
	// 消毒记录编号
	AntID string `json:"number"`
	// 内镜编号
	No string `json:"endoscope_number"`
	// 内镜类型
	Type string `json:"endoscope_type"`
	// 操作员
	Operator string `json:"operator"`
	// 病人姓名
	PatientName string `json:"patient"`
	// 医生姓名
	DoctorName string `json:"doctor"`
	// 传染病
	Diseases int `json:"diseases,string"`
	// 当前操作 步骤名称
	CurStepName string `json:"curstepname"`
	// 当前操作 步骤 时间长
	CurStepTimeLong string `json:"cstl"`
	//已经完成的操作步骤
	Ant_Steps []Ant_Step `json:"AntSteps"`
	//已经完成的操作步骤 数量
	StepCount int `json:"stepcount"`
}
type Ant_Step struct {
	Step     string `json:"Step"`
	CostTime string `json:"CostTime"`
}

func generateMonitorHash(antList map[string]runtime.AntInfo) string {
	data, _ := json.Marshal(antList)   // 将当前状态序列化为JSON
	hash := md5.Sum(data)              // 使用MD5生成哈希值
	return hex.EncodeToString(hash[:]) // 将字节切片转换为十六进制字符串
}

func Get(ginC *gin.Context) {
	//测试
	if *common.TestMode {
		generateTextData()
	}

	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var antList []Ant_Info
	for _, key := range runtime.Sequence_key {
		ant := Ant_Info{}
		runtime.G_Processing_Ant_RWMutex.RLock()
		keyAnt := runtime.G_Processing_Ant[key]
		runtime.G_Processing_Ant_RWMutex.RUnlock()

		ant.AntID = keyAnt.AntID
		ant.No = keyAnt.No
		ant.Voice_NO = keyAnt.VoiceNo
		ant.Type = keyAnt.Type
		ant.Operator = keyAnt.Operator
		ant.PatientName = keyAnt.PatientName
		if ant.PatientName == "" {
			ant.PatientName = "　"
		}
		ant.DoctorName = keyAnt.DoctorName
		if ant.DoctorName == "" {
			ant.DoctorName = "　"
		}
		ant.Diseases = keyAnt.Diseases

		ant.CurStepName = keyAnt.CurStepName
		if keyAnt.CurStepName == "消毒" && keyAnt.Diseases == 2 {
			ant.CurStepName = "消毒(灭菌)"
		}
		if strings.Contains(ant.CurStepName, "洗镜机") {
			timeLong := time.Now().Unix() - keyAnt.CurStepStartTime.Unix()
			program := runtime.G_Program_Info.Program[ant.CurStepName]
			var accumulate int
			tempStep := ""
			for x := 0; x < len(program.StepList); x++ {
				accumulate += program.StepList[strconv.Itoa(x+1)].Cost_time
				tempStep = program.StepList[strconv.Itoa(x+1)].Step
				if accumulate >= int(timeLong) {
					break
				}
			}
			ant.CurStepName = strings.Replace(ant.CurStepName, "洗镜机", "No", -1) + ":" + tempStep
		}
		ant.CurStepTimeLong = common.SecToStr(time.Now().Unix() - keyAnt.CurStepStartTime.Unix())
		ant.StepCount = 0
		for i, res := range keyAnt.StepList {
			if res == "消毒" && keyAnt.Diseases == 2 {
				res = "消毒(灭菌)"
			}
			ant.Ant_Steps = append(ant.Ant_Steps, Ant_Step{res,
				common.SecToStr(int64(keyAnt.StepTimeLongList[i]))})
			ant.StepCount++
		}
		antList = append(antList, ant)
	}

	// 不做hash 校验了
	hash := generateMonitorHash(runtime.G_Processing_Ant)
	//if post["hash"].(string) == hash {
	//	ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "与上次内容相同", "hash": hash,
	//		"antList": []Ant_Info{}, "datetime": time.Now().Format("2006-01-02 15:04:05")})
	//	return
	//}
	if len(antList) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "无新数据", "hash": hash,
			"antList": []Ant_Info{}, "datetime": time.Now().Format("2006-01-02 15:04:05")})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "有数据更新", "hash": hash,
		"antList": antList, "datetime": time.Now().Format("2006-01-02 15:04:05")})
}

func GetGeneralInformation(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功",
		"error":   common.MsgQueueError.PopOrWait(1),
		"warning": common.MsgQueueWarning.PopOrWait(1),
		"info":    common.MsgQueueInfo.PopOrWait(1)})
}

type AntModifyData struct {
	Number          string   `json:"number"`
	EndoscopeNumber string   `json:"endoscope_number"`
	PatientName     string   `json:"patient_name"`
	DocName         string   `json:"doc_name"`
	Diseases        string   `json:"diseases"`
	LeakPart        []string `json:"leakPart"`
}

func AntModify(ginC *gin.Context) {
	var post AntModifyData
	err := ginC.ShouldBind(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	if post.EndoscopeNumber != "" {
		runtime.G_Processing_Ant_RWMutex.Lock()
		ant := runtime.G_Processing_Ant[post.EndoscopeNumber]
		if ant.AntID != "" {
			ant.PatientName = post.PatientName
			ant.DoctorName = post.DocName
			ant.Diseases, _ = strconv.Atoi(post.Diseases)
			runtime.G_Processing_Ant[post.EndoscopeNumber] = ant
		}
		runtime.G_Processing_Ant_RWMutex.Unlock()
	}
	runtime.AntDataUpdate(post.Number, post.PatientName, post.DocName, post.Diseases)
	runtime.AntDataLeakPart(post.Number, strings.Join(post.LeakPart, ","))

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功"})
}

func GetLeakPart(ginC *gin.Context) {
	var post map[string]string
	err := ginC.Bind(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": runtime.GetLeakPart(post["number"])})
}

var onlyone = false

func generateRandomStepLists() ([]string, []int) {
	stepNames := []string{"测漏", "初洗", "清洗", "漂洗", "消毒", "末洗", "酒精", "干燥", "超声波"}
	var steps []string
	var times []int
	count := rand.Intn(4) + 1 // 生成1到4步之间的随机步骤数
	for i := 0; i < count; i++ {
		steps = append(steps, stepNames[rand.Intn(len(stepNames))])
		times = append(times, rand.Intn(120)) // 随机时长在0到120秒之间
	}
	return steps, times
}
func generateTextData() {
	if onlyone {
		return
	}
	onlyone = true

	runtime.G_Processing_Ant = map[string]runtime.AntInfo{}
	runtime.Sequence_key = []string{}

	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < 15; i++ {
		stepList, stepTimeLongList := generateRandomStepLists()
		record := runtime.AntInfo{
			VoiceNo:          fmt.Sprintf("V%03d", i+1),
			AntID:            "ANT2023000033",
			No:               fmt.Sprintf("E%08d", rand.Intn(100)),
			Type:             "胃镜",
			Info:             "内窥镜描述信息",
			Operator:         "操作员" + strconv.Itoa(i+1),
			PatientName:      "是病人" + strconv.Itoa(i+1),
			DoctorName:       "是医生" + strconv.Itoa(i+1),
			Diseases:         rand.Intn(3),
			CurStepName:      "消毒",
			CurStepIP:        fmt.Sprintf("192.168.1.%d", i+1),
			CurStepStartTime: time.Now(),
			CurStepTimeLong:  rand.Intn(1500),
			StepList:         stepList,
			StepTimeLongList: stepTimeLongList,
		}
		runtime.G_Processing_Ant_RWMutex.Lock()
		runtime.G_Processing_Ant[record.No] = record
		runtime.G_Processing_Ant_RWMutex.Unlock()
		//使用list 记录当前进行中的洗消记录
		runtime.Sequence_key = append(runtime.Sequence_key, record.No)
	}
}
