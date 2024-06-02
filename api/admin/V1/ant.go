package V1

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"time"
)

func GetDeviceTypeList(ginC *gin.Context) {
	res, err := models.DeviceTypelist()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
func GetOperatorList(ginC *gin.Context) {
	res, err := models.Operatorlist()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
func GetDoctorList(ginC *gin.Context) {
	res, err := models.Doctorlist()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
func GetEcmList(ginC *gin.Context) {
	res, err := models.EcmList()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

type AntListCondition struct {
	Date            []time.Time `form:"date"`
	DeviceNo        string      `form:"deviceNo"`
	DeviceType      string      `form:"deviceType"`
	DeviceInfo      string      `form:"deviceInfo"`
	Diseases        string      `form:"diseases"`
	Number          string      `form:"number"`
	Operator        string      `form:"operator"`
	Patient_name    string      `form:"patient_name"`
	Doctor          string      `form:"doctor"`
	Total_cost_time int         `form:"total_cost_time"`
	Ecm_name        string      `form:"ecm_name"`
	Query_range     string      `form:"query_range"`
	Pagination      struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryAntList(ginC *gin.Context) {
	var alc AntListCondition
	err := ginC.BindJSON(&alc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	alc.Date[0] = alc.Date[0].In(time.Local)
	alc.Date[1] = alc.Date[1].In(time.Local)

	conditionStr := " where t.begin_time between '" + alc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + alc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(alc.DeviceNo) > 0 {
		conditionStr += " and t.endoscope_number like '%" + alc.DeviceNo + "%' "
	}
	if len(alc.DeviceType) > 0 {
		conditionStr += " and t.endoscope_type = '" + alc.DeviceType + "' "
	}
	if len(alc.DeviceInfo) > 0 {
		conditionStr += " and t.endoscope_info like '%" + alc.DeviceInfo + "%' "
	}
	if len(alc.Number) > 0 {
		conditionStr += " and t.number like '%" + alc.Number + "%' "
	}
	if len(alc.Operator) > 0 {
		conditionStr += " and t.operator = '" + alc.Operator + "' "
	}
	if len(alc.Patient_name) > 0 {
		conditionStr += " and t.patient_name like '%" + alc.Patient_name + "%' "
	}
	if len(alc.Doctor) > 0 {
		conditionStr += " and t.doc_name = '" + alc.Doctor + "' "
	}
	if len(alc.Ecm_name) > 0 {
		conditionStr += " and t.number in (select DISTINCT number from ant_step where number=t.number and  washing_machine='" + alc.Ecm_name + "') "
	}
	if alc.Total_cost_time > 0 {
		conditionStr += " and t.total_cost_time <=" + strconv.Itoa(alc.Total_cost_time)
	}
	if len(alc.Diseases) > 0 && alc.Diseases != "0" {
		conditionStr += " and t.diseases =" + alc.Diseases
	}
	if alc.Query_range == "1" { //晨洗
		conditionStr += " and begin_time = (SELECT MIN(begin_time) FROM ant WHERE  endoscope_number = t.endoscope_number " +
			"AND DATE(begin_time) = DATE(t.begin_time)) "
	}

	if alc.Query_range == "2" { //终洗
		conditionStr += " and begin_time = (SELECT MAX(begin_time) FROM ant WHERE  endoscope_number = t.endoscope_number " +
			"AND DATE(begin_time) = DATE(t.begin_time)) "
	}

	var offset int
	if alc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (alc.Pagination.CurrentPage - 1) * alc.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from ant as t " + conditionStr).Count(&total)

	listSqlStr := "select id,number,endoscope_number,endoscope_type,operator,patient_name,doc_name,diseases,begin_time" +
		",end_time,SEC_TO_TIME(total_cost_time) as total_cost_time,endoscope_info," +
		"(SELECT SEC_TO_TIME(COALESCE(SUM(cost_time),0)) FROM ant_step WHERE number = t.number AND washing_machine='') as manual_time," +
		"(SELECT SEC_TO_TIME(COALESCE(SUM(cost_time),0)) FROM ant_step WHERE number = t.number AND washing_machine<>'') as machine_time " +
		"from ant as t "

	// 分页排序
	ordLimitStr := " order by " + alc.Sort + " " + alc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(alc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": alc.Pagination.CurrentPage, "pageSize": alc.Pagination.PageSize, "antList": rows}})
}

func QueryAntSteps(ginC *gin.Context) {
	var res []map[string]interface{}
	err := ginC.Bind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var steps []models.AntStep
	db := orm.Where("number = ?", res[0]["antNumber"]).Order("id ").Find(&steps)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"steps": &steps}})
}

func AntModify(ginC *gin.Context) {
	var res []models.Ant
	err := ginC.Bind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	db := orm.Model(&models.Ant{}).Select("PatientName", "DocName", "Diseases").Where("Number = ?", res[0].Number).Updates(res[0])
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type ABE struct {
	Alc     AntListCondition `form:"alc"`
	ABEType int              `form:"abeType"`
}

func AntBatchExport(ginC *gin.Context) {
	var alc ABE
	err := ginC.Bind(&alc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	alc.Alc.Date[0] = alc.Alc.Date[0].In(time.Local)
	alc.Alc.Date[1] = alc.Alc.Date[1].In(time.Local)

	conditionStr := " where t.begin_time between '" + alc.Alc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + alc.Alc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(alc.Alc.DeviceNo) > 0 {
		conditionStr += " and t.endoscope_number like '%" + alc.Alc.DeviceNo + "%' "
	}
	if len(alc.Alc.DeviceType) > 0 {
		conditionStr += " and t.endoscope_type = '" + alc.Alc.DeviceType + "' "
	}
	if len(alc.Alc.DeviceInfo) > 0 {
		conditionStr += " and t.endoscope_info like '%" + alc.Alc.DeviceInfo + "%' "
	}
	if len(alc.Alc.Number) > 0 {
		conditionStr += " and t.number like '%" + alc.Alc.Number + "%' "
	}
	if len(alc.Alc.Operator) > 0 {
		conditionStr += " and t.operator = '" + alc.Alc.Operator + "' "
	}
	if len(alc.Alc.Patient_name) > 0 {
		conditionStr += " and t.patient_name like '%" + alc.Alc.Patient_name + "%' "
	}
	if len(alc.Alc.Doctor) > 0 {
		conditionStr += " and t.doc_name = '" + alc.Alc.Doctor + "' "
	}
	if len(alc.Alc.Ecm_name) > 0 {
		conditionStr += " and t.number in (select DISTINCT number from ant_step where number=t.number and  washing_machine='" + alc.Alc.Ecm_name + "') "
	}
	if alc.Alc.Total_cost_time > 0 {
		conditionStr += " and t.total_cost_time <=" + strconv.Itoa(alc.Alc.Total_cost_time)
	}
	if len(alc.Alc.Diseases) > 0 && alc.Alc.Diseases != "0" {
		conditionStr += " and t.diseases =" + alc.Alc.Diseases
	}
	if alc.Alc.Query_range == "1" { //晨洗
		conditionStr += " and begin_time = (SELECT MIN(begin_time) FROM ant WHERE  endoscope_number = t.endoscope_number " +
			"AND DATE(begin_time) = DATE(t.begin_time)) "
	}

	if alc.Alc.Query_range == "2" { //终洗
		conditionStr += " and begin_time = (SELECT MAX(begin_time) FROM ant WHERE  endoscope_number = t.endoscope_number " +
			"AND DATE(begin_time) = DATE(t.begin_time)) "
	}

	listSqlStr := "select id,number,endoscope_number,endoscope_type,operator,patient_name,doc_name,diseases,begin_time" +
		",end_time,SEC_TO_TIME(total_cost_time) as total_cost_time,endoscope_info," +
		"(SELECT SEC_TO_TIME(COALESCE(SUM(cost_time),0)) FROM ant_step WHERE number = t.number AND washing_machine='') as manual_time," +
		"(SELECT SEC_TO_TIME(COALESCE(SUM(cost_time),0)) FROM ant_step WHERE number = t.number AND washing_machine<>'') as machine_time " +
		"from ant as t "

	// 排序
	ordLimitStr := " order by " + alc.Alc.Sort + " " + alc.Alc.Order

	var rows []map[string]interface{}
	orm, _ := models.NewOrm()
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)

	if db.Error != nil || len(rows) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "没有数据"})
		return
	}
	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/AntExport2*"})
	excelFileName := "./bgStatic/AntExport" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/ExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "没有数据"})
		return
	}
	_ = file.SetCellValue("洗消数据", "C4", time.Now())
	_ = file.SetCellValue("洗消数据", "E4", len(rows))
	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("洗消数据", "A"+strconv.Itoa(rowindex), val["id"].(int64))
		_ = file.SetCellValue("洗消数据", "B"+strconv.Itoa(rowindex), val["number"].(string))
		_ = file.SetCellValue("洗消数据", "C"+strconv.Itoa(rowindex), val["endoscope_number"].(string))
		_ = file.SetCellValue("洗消数据", "D"+strconv.Itoa(rowindex), val["endoscope_type"].(string))
		_ = file.SetCellValue("洗消数据", "E"+strconv.Itoa(rowindex), val["endoscope_info"].(string))
		_ = file.SetCellValue("洗消数据", "F"+strconv.Itoa(rowindex), val["operator"].(string))
		_ = file.SetCellValue("洗消数据", "G"+strconv.Itoa(rowindex), val["patient_name"].(string))
		_ = file.SetCellValue("洗消数据", "H"+strconv.Itoa(rowindex), val["doc_name"].(string))
		if val["diseases"].(int64) == 2 {
			_ = file.SetCellValue("洗消数据", "I"+strconv.Itoa(rowindex), "是")
		} else if val["diseases"].(int64) == 1 {
			_ = file.SetCellValue("洗消数据", "I"+strconv.Itoa(rowindex), "否")
		} else {
			_ = file.SetCellValue("洗消数据", "I"+strconv.Itoa(rowindex), "未知")
		}
		_ = file.SetCellValue("洗消数据", "J"+strconv.Itoa(rowindex), val["begin_time"].(time.Time))
		_ = file.SetCellValue("洗消数据", "K"+strconv.Itoa(rowindex), val["manual_time"].(string))
		_ = file.SetCellValue("洗消数据", "L"+strconv.Itoa(rowindex), val["machine_time"].(string))
		_ = file.SetCellValue("洗消数据", "M"+strconv.Itoa(rowindex), val["end_time"].(time.Time))
		_ = file.SetCellValue("洗消数据", "N"+strconv.Itoa(rowindex), val["total_cost_time"].(string))

		rowindex++

		if alc.ABEType == 1 {
			var steps []models.AntStep
			orm.Where("number = ?", val["number"].(string)).Order("id ").Find(&steps)
			for _, step_val := range steps {
				_ = file.SetCellValue("洗消数据", "C"+strconv.Itoa(rowindex), "步骤：")
				_ = file.SetCellValue("洗消数据", "D"+strconv.Itoa(rowindex), step_val.Step)
				_ = file.SetCellValue("洗消数据", "E"+strconv.Itoa(rowindex), "时长：")
				_ = file.SetCellValue("洗消数据", "F"+strconv.Itoa(rowindex), step_val.CostTime)
				_ = file.SetCellValue("洗消数据", "G"+strconv.Itoa(rowindex), step_val.WashingMachine)
				rowindex++
			}
		}
		if alc.ABEType == 2 {
			var steps []models.AntStep
			orm.Where("number = ?", val["number"].(string)).Order("id ").Find(&steps)
			for _, step_val := range steps {
				step_str := step_val.Step + ":" + common.SecToStr(int64(step_val.CostTime)) +
					step_val.WashingMachine

				_ = file.SetCellValue("洗消数据", "B"+strconv.Itoa(rowindex), step_str)
				rowindex++
			}
		}
	}
	_ = file.Save()
	//ginC.Header("Content-Type", "application/vnd.ms-excel")
	//ginC.Header("Content-Disposition", fmt.Sprintf("attachment; filename=AntExport%s.xlsx", time.Now().Format("2006-01-02")))
	//ginC.File(excelFileName)
	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/AntExport" + time.Now().Format("2006-01-02") + ".xlsx"})
}
