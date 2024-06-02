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

type LeakDetectionCondition struct {
	Date       []time.Time `json:"date"`
	DeviceNo   string      `json:"deviceNo"`
	DeviceType string      `json:"deviceType"`
	DeviceInfo string      `json:"deviceInfo"`
	Status     string      `json:"status"`
	Parts      string      `json:"parts"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryLeakDetectionList(ginC *gin.Context) {
	var ldc LeakDetectionCondition
	err := ginC.ShouldBind(&ldc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	ldc.Date[0] = ldc.Date[0].In(time.Local)
	ldc.Date[1] = ldc.Date[1].In(time.Local)

	var offset int
	if ldc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (ldc.Pagination.CurrentPage - 1) * ldc.Pagination.PageSize
	}
	ordLimitStr := " order by " + ldc.Sort + " " + ldc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(ldc.Pagination.PageSize), 10)

	conditionStr := " where l.id>0 and a.begin_time between '" + ldc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + ldc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(ldc.DeviceNo) > 0 {
		conditionStr += " and a.endoscope_number like '%" + ldc.DeviceNo + "%' "
	}
	if len(ldc.DeviceType) > 0 {
		conditionStr += " and a.endoscope_type = '" + ldc.DeviceType + "' "
	}
	if len(ldc.DeviceInfo) > 0 {
		conditionStr += " and a.endoscope_info like '%" + ldc.DeviceInfo + "%' "
	}
	if len(ldc.Status) > 0 {
		conditionStr += " and l.conclusion = '" + ldc.Status + "' "
	}
	if len(ldc.Parts) > 0 {
		conditionStr += " and l.leakparts like '%" + ldc.Parts + "%' "
	}
	listSqlStr := "select a.id from ant as a  left join leak_detection_record as l on a.number=l.number "
	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw(listSqlStr + conditionStr).Count(&total)

	listSqlStr = "select a.id,a.number,l.id as lid,endoscope_number,endoscope_type,endoscope_info,begin_time,operator,conclusion,leakparts from ant as a left join leak_detection_record as l on a.number=l.number "
	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "data": rows,
		"Pagination": gin.H{"total": total, "currentPage": ldc.Pagination.CurrentPage, "pageSize": ldc.Pagination.PageSize}})
}

func ModifyPartsModify(ginC *gin.Context) {
	var upi models.LeakDetectionRecord
	err := ginC.ShouldBind(&upi)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}

	orm.Save(upi)
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

func LeakDetectionExport(ginC *gin.Context) {
	var ldc LeakDetectionCondition
	err := ginC.ShouldBind(&ldc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	ldc.Date[0] = ldc.Date[0].In(time.Local)
	ldc.Date[1] = ldc.Date[1].In(time.Local)

	listSqlStr := "select a.id,a.number,l.id as lid,endoscope_number,endoscope_type,endoscope_info,begin_time,operator,conclusion,leakparts from ant as a left join leak_detection_record as l on a.number=l.number "
	conditionStr := " where l.id>0 and a.begin_time between '" + ldc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + ldc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "
	ordLimitStr := " order by " + ldc.Sort + " " + ldc.Order

	if len(ldc.DeviceNo) > 0 {
		conditionStr += " and a.endoscope_number like '%" + ldc.DeviceNo + "%' "
	}
	if len(ldc.DeviceType) > 0 {
		conditionStr += " and a.endoscope_type = '" + ldc.DeviceType + "' "
	}
	if len(ldc.DeviceInfo) > 0 {
		conditionStr += " and a.endoscope_info like '%" + ldc.DeviceInfo + "%' "
	}
	if len(ldc.Status) > 0 {
		conditionStr += " and l.conclusion = '" + ldc.Status + "' "
	}
	if len(ldc.Parts) > 0 {
		conditionStr += " and l.leakparts like '%" + ldc.Parts + "%' "
	}

	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error})
		return
	}

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/LeakExportTemplates2*"})
	excelFileName := "./bgStatic/LeakExportTemplates" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/LeakExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "生成文件错误"})
		return
	}
	_ = file.SetCellValue("测漏记录", "C4", time.Now())
	_ = file.SetCellValue("测漏记录", "E4", len(rows))

	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("测漏记录", "A"+strconv.Itoa(rowindex), val["id"])
		_ = file.SetCellValue("测漏记录", "B"+strconv.Itoa(rowindex), val["number"])
		_ = file.SetCellValue("测漏记录", "C"+strconv.Itoa(rowindex), val["endoscope_number"])
		_ = file.SetCellValue("测漏记录", "D"+strconv.Itoa(rowindex), val["endoscope_type"])
		_ = file.SetCellValue("测漏记录", "E"+strconv.Itoa(rowindex), val["endoscope_info"])
		_ = file.SetCellValue("测漏记录", "F"+strconv.Itoa(rowindex), val["operator"])
		_ = file.SetCellValue("测漏记录", "G"+strconv.Itoa(rowindex), val["begin_time"])
		_ = file.SetCellValue("测漏记录", "H"+strconv.Itoa(rowindex), val["conclusion"])
		_ = file.SetCellValue("测漏记录", "I"+strconv.Itoa(rowindex), val["leakparts"])
		rowindex++
	}
	_ = file.Save()

	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/LeakExportTemplates" +
		time.Now().Format("2006-01-02") + ".xlsx"})
}

type StorageCondition struct {
	DeviceNo   string `json:"deviceNo"`
	DeviceType string `json:"deviceType"`
	DeviceInfo string `json:"deviceInfo"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func GetStorageList(ginC *gin.Context) {
	var sc StorageCondition
	err := ginC.ShouldBind(&sc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	var offset int
	if sc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (sc.Pagination.CurrentPage - 1) * sc.Pagination.PageSize
	}
	ordLimitStr := " order by " + sc.Sort + " " + sc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(sc.Pagination.PageSize), 10)

	conditionStr := " where d.id>0 "

	if len(sc.DeviceNo) > 0 {
		conditionStr += " and d.endoscope_number like '%" + sc.DeviceNo + "%' "
	}
	if len(sc.DeviceType) > 0 {
		conditionStr += " and d.endoscope_type = '" + sc.DeviceType + "' "
	}
	if len(sc.DeviceInfo) > 0 {
		conditionStr += " and d.endoscope_info like '%" + sc.DeviceInfo + "%' "
	}

	listSqlStr := "select id from device_info as d "
	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw(listSqlStr + conditionStr).Count(&total)

	listSqlStr = "SELECT d.id,d.endoscope_number,d.endoscope_type,d.endoscope_info," +
		"IFNULL((SELECT access from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as status1," +
		"IFNULL((SELECT operator from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as operator1," +
		"IFNULL((SELECT time from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as time1," +
		"IFNULL((SELECT access from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as status2," +
		"IFNULL((SELECT operator from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as operator2," +
		"IFNULL((SELECT time from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as time2 " +
		"FROM device_info AS d "

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "data": rows,
		"Pagination": gin.H{"total": total, "currentPage": sc.Pagination.CurrentPage, "pageSize": sc.Pagination.PageSize}})
}

type StorageRecordsCondition struct {
	Number     string      `json:"number"`
	Date       []time.Time `json:"date"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func GetStorageRecordList(ginC *gin.Context) {
	var src StorageRecordsCondition
	err := ginC.ShouldBind(&src)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	src.Date[0] = src.Date[0].In(time.Local)
	src.Date[1] = src.Date[1].In(time.Local)

	var offset int
	if src.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (src.Pagination.CurrentPage - 1) * src.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Model(&models.StoreLog{}).Select("id").
		Where("time BETWEEN ? AND ? AND endoscope_number=?", src.Date[0], src.Date[1], src.Number).Count(&total)

	var rows []models.StoreLog
	orm.Where("time BETWEEN ? AND ? AND endoscope_number=?", src.Date[0], src.Date[1], src.Number).
		Order(src.Sort + " " + src.Order).Offset(offset).Limit(src.Pagination.PageSize).Find(&rows)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "data": rows, "Pagination": gin.H{"total": total,
		"currentPage": src.Pagination.CurrentPage, "pageSize": src.Pagination.PageSize}})
}

func StorageExport(ginC *gin.Context) {
	var sc StorageCondition
	err := ginC.ShouldBind(&sc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	ordLimitStr := " order by " + sc.Sort + " " + sc.Order
	conditionStr := " where d.id>0 "

	if len(sc.DeviceNo) > 0 {
		conditionStr += " and d.endoscope_number like '%" + sc.DeviceNo + "%' "
	}
	if len(sc.DeviceType) > 0 {
		conditionStr += " and d.endoscope_type = '" + sc.DeviceType + "' "
	}
	if len(sc.DeviceInfo) > 0 {
		conditionStr += " and d.endoscope_info like '%" + sc.DeviceInfo + "%' "
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	listSqlStr := "SELECT d.id,d.endoscope_number,d.endoscope_type,d.endoscope_info," +
		"IFNULL((SELECT access from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as status1," +
		"IFNULL((SELECT operator from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as operator1," +
		"IFNULL((SELECT time from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 0,1),'') as time1," +
		"IFNULL((SELECT access from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as status2," +
		"IFNULL((SELECT operator from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as operator2," +
		"IFNULL((SELECT time from store_log WHERE endoscope_number=d.endoscope_number ORDER BY id Desc LIMIT 1,1),'') as time2 " +
		"FROM device_info AS d "

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/StoreTemplates2*"})
	excelFileName := "./bgStatic/StoreTemplates" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/StoreTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "生成文件错误"})
		return
	}
	_ = file.SetCellValue("出入库记录", "C4", time.Now())
	_ = file.SetCellValue("出入库记录", "E4", len(rows))

	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("出入库记录", "A"+strconv.Itoa(rowindex), val["id"])
		_ = file.SetCellValue("出入库记录", "B"+strconv.Itoa(rowindex), val["endoscope_number"])
		_ = file.SetCellValue("出入库记录", "C"+strconv.Itoa(rowindex), val["endoscope_type"])
		_ = file.SetCellValue("出入库记录", "D"+strconv.Itoa(rowindex), val["endoscope_info"])
		_ = file.SetCellValue("出入库记录", "E"+strconv.Itoa(rowindex), val["status1"])
		_ = file.SetCellValue("出入库记录", "F"+strconv.Itoa(rowindex), val["operator1"])
		_ = file.SetCellValue("出入库记录", "G"+strconv.Itoa(rowindex), val["time1"])
		_ = file.SetCellValue("出入库记录", "H"+strconv.Itoa(rowindex), val["status2"])
		_ = file.SetCellValue("出入库记录", "I"+strconv.Itoa(rowindex), val["operator2"])
		_ = file.SetCellValue("出入库记录", "J"+strconv.Itoa(rowindex), val["time2"])
		rowindex++
	}
	_ = file.Save()

	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/StoreTemplates" +
		time.Now().Format("2006-01-02") + ".xlsx"})
}

func RepairRecordExport(ginC *gin.Context) {
	var dlc DeviceListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(dlc.DeviceNo) > 0 {
		conditionStr += " and endoscope_number like '%" + dlc.DeviceNo + "%' "
	}
	if len(dlc.DeviceType) > 0 {
		conditionStr += " and endoscope_type = '" + dlc.DeviceType + "' "
	}
	if len(dlc.DeviceInfo) > 0 {
		conditionStr += " and endoscope_info like '%" + dlc.DeviceInfo + "%' "
	}
	if dlc.Status != "0" {
		conditionStr += " and status =" + dlc.Status
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	listSqlStr := "select * from device_info "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/DeviceExportTemplates2*"})
	excelFileName := "./bgStatic/DeviceExportTemplates" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/DeviceExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "生成文件错误"})
		return
	}
	_ = file.SetCellValue("内窥镜记录", "C4", time.Now())
	_ = file.SetCellValue("内窥镜记录", "E4", len(rows))

	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("内窥镜记录", "A"+strconv.Itoa(rowindex), val["id"])
		_ = file.SetCellValue("内窥镜记录", "B"+strconv.Itoa(rowindex), val["endoscope_number"])
		_ = file.SetCellValue("内窥镜记录", "C"+strconv.Itoa(rowindex), val["endoscope_type"])
		_ = file.SetCellValue("内窥镜记录", "D"+strconv.Itoa(rowindex), val["endoscope_info"])
		status := "正常"
		if val["status"].(int32) == 2 {
			status = "维修"
		}
		_ = file.SetCellValue("内窥镜记录", "E"+strconv.Itoa(rowindex), status)
		rowindex++
	}
	_ = file.Save()

	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/DeviceExportTemplates" +
		time.Now().Format("2006-01-02") + ".xlsx"})
}

func SendRepairModify(ginC *gin.Context) {
	var rr models.RepairRecords
	err := ginC.ShouldBind(&rr)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}

	db := orm.Save(&rr)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	orm.Model(&models.DeviceInfo{}).Where("endoscope_number = ?", rr.EndoscopeNumber).Update("status", 2)
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "修改成功"})
}

func CompletedRepairModify(ginC *gin.Context) {
	var rr models.RepairRecords
	err := ginC.ShouldBind(&rr)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}

	db := orm.Model(&models.RepairRecords{}).Where("endoscope_number = ? AND return_date is null", rr.EndoscopeNumber).
		Order("id desc").Limit(1).Updates(models.RepairRecords{ReturnDate: rr.ReturnDate, RepairCost: rr.RepairCost, RepairNotes: rr.RepairNotes})
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	orm.Model(&models.DeviceInfo{}).Where("endoscope_number = ?", rr.EndoscopeNumber).Update("status", 1)
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "修改成功"})
}

type RepairDetailsCondition struct {
	Number     string      `json:"endoscope_number"`
	Date       []time.Time `json:"date"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryRepairDetailsList(ginC *gin.Context) {
	var src RepairDetailsCondition
	err := ginC.ShouldBind(&src)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	src.Date[0] = src.Date[0].In(time.Local)
	src.Date[1] = src.Date[1].In(time.Local)
	var offset int
	if src.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (src.Pagination.CurrentPage - 1) * src.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Model(&models.RepairRecords{}).Select("id").
		Where("send_out_date BETWEEN ? AND ? AND endoscope_number=?", src.Date[0], src.Date[1], src.Number).Count(&total)

	var rows []models.RepairRecords
	orm.Where("send_out_date BETWEEN ? AND ? AND endoscope_number=?", src.Date[0], src.Date[1], src.Number).
		Order(src.Sort + " " + src.Order).Offset(offset).Limit(src.Pagination.PageSize).Find(&rows)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "data": rows, "Pagination": gin.H{"total": total,
		"currentPage": src.Pagination.CurrentPage, "pageSize": src.Pagination.PageSize}})
}

type RepairRecordCondition struct {
	Date            []time.Time `json:"date"`
	EndoscopeNumber string      `json:"endoscope_number"`
	EndoscopeType   string      `json:"endoscope_type"`
	EndoscopeInfo   string      `json:"endoscope_info"`
	Status          string      `json:"status"`
	ProblemLocation string      `json:"problem_location"`
	Pagination      struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryRepairRecordList(ginC *gin.Context) {
	var rrc RepairRecordCondition
	err := ginC.ShouldBind(&rrc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	rrc.Date[0] = rrc.Date[0].In(time.Local)
	rrc.Date[1] = rrc.Date[1].In(time.Local)
	var offset int
	if rrc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (rrc.Pagination.CurrentPage - 1) * rrc.Pagination.PageSize
	}

	ordLimitStr := " order by " + rrc.Sort + " " + rrc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(rrc.Pagination.PageSize), 10)
	conditionStr := " where r.send_out_date between '" + rrc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + rrc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(rrc.EndoscopeNumber) > 0 {
		conditionStr += " and d.endoscope_number like '%" + rrc.EndoscopeNumber + "%' "
	}
	if len(rrc.EndoscopeType) > 0 {
		conditionStr += " and d.endoscope_type = '" + rrc.EndoscopeType + "' "
	}
	if len(rrc.EndoscopeInfo) > 0 {
		conditionStr += " and d.endoscope_info like '%" + rrc.EndoscopeInfo + "%' "
	}
	if rrc.Status == "在修" {
		conditionStr += " and r.return_date IS NULL "
	} else if rrc.Status == "完修" {
		conditionStr += " and r.return_date IS NOT NULL "
	}
	if len(rrc.ProblemLocation) > 0 {
		conditionStr += " and r.problem_location like '%" + rrc.ProblemLocation + "%' "
	}

	listSqlStr := "select r.id from repair_records as r  left join device_info as d on r.endoscope_number=d.endoscope_number "
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}

	var total int64
	db := orm.Raw(listSqlStr + conditionStr).Count(&total)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	listSqlStr = "select r.*,d.endoscope_type,d.endoscope_info from repair_records as r left join device_info as d on r.endoscope_number=d.endoscope_number "
	var rows []map[string]interface{}
	db = orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "data": rows, "Pagination": gin.H{"total": total,
		"currentPage": rrc.Pagination.CurrentPage, "pageSize": rrc.Pagination.PageSize}})
}

func RepairManagementExport(ginC *gin.Context) {
	var rrc RepairRecordCondition
	err := ginC.ShouldBind(&rrc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	rrc.Date[0] = rrc.Date[0].In(time.Local)
	rrc.Date[1] = rrc.Date[1].In(time.Local)

	ordLimitStr := " order by " + rrc.Sort + " " + rrc.Order
	conditionStr := " where r.send_out_date between '" + rrc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + rrc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(rrc.EndoscopeNumber) > 0 {
		conditionStr += " and d.endoscope_number like '%" + rrc.EndoscopeNumber + "%' "
	}
	if len(rrc.EndoscopeType) > 0 {
		conditionStr += " and d.endoscope_type = '" + rrc.EndoscopeType + "' "
	}
	if len(rrc.EndoscopeInfo) > 0 {
		conditionStr += " and d.endoscope_info like '%" + rrc.EndoscopeInfo + "%' "
	}
	if rrc.Status == "在修" {
		conditionStr += " and r.return_date IS NULL "
	} else if rrc.Status == "完修" {
		conditionStr += " and r.return_date IS NOT NULL "
	}
	if len(rrc.ProblemLocation) > 0 {
		conditionStr += " and r.problem_location like '%" + rrc.ProblemLocation + "%' "
	}

	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	listSqlStr := "select r.*,d.endoscope_type,d.endoscope_info from repair_records as r left join device_info as d on r.endoscope_number=d.endoscope_number "
	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/RepairExportTemplates2*"})
	excelFileName := "./bgStatic/RepairExportTemplates" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/RepairExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "生成文件错误"})
		return
	}
	_ = file.SetCellValue("维修记录", "C4", time.Now())
	_ = file.SetCellValue("维修记录", "E4", len(rows))

	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("维修记录", "A"+strconv.Itoa(rowindex), val["id"])
		_ = file.SetCellValue("维修记录", "B"+strconv.Itoa(rowindex), val["endoscope_number"])
		_ = file.SetCellValue("维修记录", "C"+strconv.Itoa(rowindex), val["endoscope_type"])
		_ = file.SetCellValue("维修记录", "D"+strconv.Itoa(rowindex), val["endoscope_info"])
		if val["return_date"] == nil {
			_ = file.SetCellValue("维修记录", "E"+strconv.Itoa(rowindex), "在修")
		} else {
			_ = file.SetCellValue("维修记录", "E"+strconv.Itoa(rowindex), "完修")
		}
		_ = file.SetCellValue("维修记录", "F"+strconv.Itoa(rowindex), val["responsible_person"])
		_ = file.SetCellValue("维修记录", "G"+strconv.Itoa(rowindex), val["problem_description"])
		_ = file.SetCellValue("维修记录", "H"+strconv.Itoa(rowindex), val["problem_location"])
		if val["return_date"] == nil {
			_ = file.SetCellValue("维修记录", "I"+strconv.Itoa(rowindex), "-")
		} else {
			difference := val["return_date"].(time.Time).Sub(val["send_out_date"].(time.Time))
			_ = file.SetCellValue("维修记录", "I"+strconv.Itoa(rowindex), int64(difference.Hours()/24))
		}
		_ = file.SetCellValue("维修记录", "J"+strconv.Itoa(rowindex), val["send_out_date"])
		_ = file.SetCellValue("维修记录", "K"+strconv.Itoa(rowindex), val["return_date"])
		_ = file.SetCellValue("维修记录", "L"+strconv.Itoa(rowindex), val["repair_cost"])
		_ = file.SetCellValue("维修记录", "M"+strconv.Itoa(rowindex), val["repair_notes"])
		rowindex++
	}
	_ = file.Save()

	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/RepairExportTemplates" +
		time.Now().Format("2006-01-02") + ".xlsx"})
}
