package runtime

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"time"
)

type AntDb struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`

	A_steps map[string]AntStepDb `json:"steps"`
}

// AntDb_1 2021-10-7 返回数据增加 endoscope_model
type AntDb_1 struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`
	A_endoscope_model  string `json:"endoscope_model"`

	A_steps map[string]AntStepDb `json:"steps"`
}

type AntDbSlic struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`

	A_steps []AntStepDb `json:"steps"`
}

type AntDbSlic_1 struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`
	A_endoscope_model  string `json:"endoscope_model"`

	A_steps []AntStepDb `json:"steps"`
}

type AntDbSlic_2 struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`
	A_endoscope_model  string `json:"endoscope_model"`

	A_patientVisitId string `json:"patientVisitId"`
	A_checkReportId  string `json:"checkReportId"`

	A_steps []AntStepDb `json:"steps"`
}

type AntStepDb struct {
	S_id              string `json:"step_id"`
	S_number          string `json:"step_number"`
	S_step            string `json:"step"`
	S_stepip          string `json:"stepip"`
	S_cost_time       string `json:"step_cost_time"`
	S_washing_machine string `json:"step_washing_machine"`
}

type AntInfo struct {
	// 语音播报编号
	VoiceNo string `json:"-"`
	// 消毒记录编号
	AntID string `json:"number"`
	// 内镜编号
	No string `json:"endoscope_number"`
	// 内镜类型
	Type string `json:"endoscope_type"`
	// 内镜信息
	Info string `json:"endoscope_info"`
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
	// 当前操作 步骤IP
	CurStepIP string `json:"curstepip"`
	// 当前操作 步骤 开始时间
	CurStepStartTime time.Time `json:"csst"`
	// 当前操作 步骤 限制结束时间  如不限制结束，设置为 0 以秒为单位。
	CurStepTimeLong int `json:"cstl"`
	// 当前内窥镜 已经完成的操作步骤 名称列表
	StepList []string `json:"StepList"`
	// 当前内窥镜 已经完成的操作步骤 时长列表
	StepTimeLongList []int `json:"step_timelonglist"`
}

func (ai *AntInfo) Init(nEndoscopyNO string, ncur_Operator string,
	neType string, neInfo string) (res bool) {

	ai.AntID = ai.GenerateAntNumber()
	if ai.AntID == "err" {
		return false
	}

	ai.No = nEndoscopyNO
	ai.Type = neType
	ai.Info = neInfo
	ai.Operator = ncur_Operator
	ai.VoiceNo = common.GetDeviceVoice(nEndoscopyNO, neInfo)
	ai.PatientName = ""
	ai.DoctorName = ""
	ai.Diseases = 0
	recordTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "insert into ant(number, endoscope_number, endoscope_type, endoscope_info," +
		" operator,patient_name,doc_name,diseases, begin_time, end_time,total_cost_time) values (" +
		"'" + ai.AntID + "'," +
		"'" + ai.No + "'," +
		"'" + ai.Type + "'," +
		"'" + ai.Info + "'," +
		"'" + ai.Operator + "','','',0," +
		"'" + recordTime + "'," +
		"'" + recordTime + "',0)"

	engine, _ := models.NewOrm()
	engine.Exec(sqlStr)

	return true
}

// InitLeakDetectionRecord 生成洗消编号
func (ai *AntInfo) InitLeakDetectionRecord(number string) {

	sql_str := "insert into leak_detection_record(number, conclusion, leakparts) values ('" + number + "','完好','')"

	engine, _ := models.NewOrm()
	engine.Exec(sql_str)
}

// GenerateAntNumber 生成洗消编号
func (ai *AntInfo) GenerateAntNumber() string {

	if time.Now().Year() < 2017 {
		return "err"
	}
	HostIdentity := common.Config.GetString("Equipment.Identity")

	var ant models.Ant
	engine, _ := models.NewOrm()
	engine.Where("number like ?", "%"+HostIdentity+"%").Last(&ant)
	fmt.Println(ant)

	if ant.Id <= 0 {
		zap.S().Errorln("数据库可能读取失败，导致序号从1开始")
		return "err"
	}

	if ant.Id > 0 {
		number := ant.Number
		reg := regexp.MustCompile(`[^` + HostIdentity + `]{4}`) // ANT2015000001 2015
		nYear := reg.FindString(number)
		reg = regexp.MustCompile(`[^` + HostIdentity + `]{6}$`) // ANT2015000001 000001
		no := reg.FindString(number)
		tempYear, _ := strconv.Atoi(nYear)
		if time.Now().Year() <= tempYear {
			index, _ := strconv.Atoi(no)
			index++
			return HostIdentity + nYear +
				common.Substring("000000"+strconv.Itoa(index), len(strconv.Itoa(index)), 6)
		}
	}

	return HostIdentity + strconv.Itoa(time.Now().Year()) + "000001"
}

func (ai *AntInfo) Destroy() {
	//  整个洗消过程用时合计
	var totalCostTime int
	for _, val := range ai.StepTimeLongList {
		totalCostTime += val
	}

	//// 洗消结束后硬加两个步骤的 时间 34秒 2023-4-17 ！！！！
	//sqlStr := "update ant set end_time='" + time.Now().Add(time.Second*34).Format("2006-01-02 15:04:05") + "',total_cost_time='" +
	//	strconv.Itoa(totalCostTime+34) + "' where number='" + ai.AntID + "'"

	sqlStr := "update ant set end_time='" + time.Now().Format("2006-01-02 15:04:05") + "',total_cost_time='" +
		strconv.Itoa(totalCostTime) + "' where number='" + ai.AntID + "'"

	engine, _ := models.NewOrm()
	engine.Exec(sqlStr)
}
