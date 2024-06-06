package V1

import (
	"TmdtServer/common"
	"TmdtServer/models"
	"TmdtServer/runtime"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"time"
)

type DataListCondition struct {
	Date       []time.Time `json:"date"`
	DeviceCode string      `json:"device_code"`
	BedCode    string      `json:"bed_code"`
	Operator   string      `json:"operator"`
	Patient    string      `json:"patient"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryDataList(ginC *gin.Context) {
	var dlc DataListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	dlc.Date[0] = dlc.Date[0].In(time.Local)
	dlc.Date[1] = dlc.Date[1].In(time.Local)

	conditionStr := " where record_time between '" + dlc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + dlc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(dlc.DeviceCode) > 0 {
		conditionStr += " and device_code like '%" + dlc.DeviceCode + "%' "
	}
	if len(dlc.BedCode) > 0 {
		conditionStr += " and bed_code like '%" + dlc.BedCode + "%' "
	}
	if len(dlc.Operator) > 0 {
		conditionStr += " and operator like '%" + dlc.Operator + "%' "
	}
	if len(dlc.Patient) > 0 {
		conditionStr += " and patient like '%" + dlc.Patient + "%' "
	}

	var offset int
	if dlc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (dlc.Pagination.CurrentPage - 1) * dlc.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from temperature_record " + conditionStr).Count(&total)

	listSqlStr := "select * from temperature_record "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "DataList": rows}})
}

func ExportDataBatch(ginC *gin.Context) {
	var dlc DataListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	dlc.Date[0] = dlc.Date[0].In(time.Local)
	dlc.Date[1] = dlc.Date[1].In(time.Local)

	conditionStr := " where record_time between '" + dlc.Date[0].Local().Format("2006-01-02") + " 00:00:00' and '" + dlc.Date[1].Local().Format("2006-01-02") + " 23:59:59' "

	if len(dlc.DeviceCode) > 0 {
		conditionStr += " and device_code like '%" + dlc.DeviceCode + "%' "
	}
	if len(dlc.BedCode) > 0 {
		conditionStr += " and bed_code like '%" + dlc.BedCode + "%' "
	}
	if len(dlc.Operator) > 0 {
		conditionStr += " and operator like '%" + dlc.Operator + "%' "
	}
	if len(dlc.Patient) > 0 {
		conditionStr += " and patient like '%" + dlc.Patient + "%' "
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from temperature_record " + conditionStr).Count(&total)

	listSqlStr := "select * from temperature_record "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/DataExport2*"})
	excelFileName := "./bgStatic/DataExport" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/ExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "没有数据"})
		return
	}
	_ = file.SetCellValue("体温数据", "C4", time.Now())
	_ = file.SetCellValue("体温数据", "E4", len(rows))
	rowindex := 6
	for _, val := range rows {
		_ = file.SetCellValue("体温数据", "A"+strconv.Itoa(rowindex), val["id"].(int32))
		_ = file.SetCellValue("体温数据", "B"+strconv.Itoa(rowindex), val["device_code"].(string))
		_ = file.SetCellValue("体温数据", "C"+strconv.Itoa(rowindex), val["bed_code"].(string))
		_ = file.SetCellValue("体温数据", "D"+strconv.Itoa(rowindex), val["operator"].(string))
		_ = file.SetCellValue("体温数据", "E"+strconv.Itoa(rowindex), val["patient"].(string))
		_ = file.SetCellValue("体温数据", "F"+strconv.Itoa(rowindex), val["temperature1"].(float32))
		_ = file.SetCellValue("体温数据", "G"+strconv.Itoa(rowindex), val["temperature2"].(float64))
		_ = file.SetCellValue("体温数据", "H"+strconv.Itoa(rowindex), val["temperature3"].(float64))
		_ = file.SetCellValue("体温数据", "I"+strconv.Itoa(rowindex), val["record_time"].(time.Time))

		rowindex++
	}
	_ = file.Save()
	//ginC.Header("Content-Type", "application/vnd.ms-excel")
	//ginC.Header("Content-Disposition", fmt.Sprintf("attachment; filename=AntExport%s.xlsx", time.Now().Format("2006-01-02")))
	//ginC.File(excelFileName)
	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/DataExport" + time.Now().Format("2006-01-02") + ".xlsx"})
}

func GetOperatorList(ginC *gin.Context) {
	res, err := models.Operatorlist()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

func GetBedsList(ginC *gin.Context) {
	rows, _ := runtime.G_BedService.GetAllBeds()
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}
