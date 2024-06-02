package V1

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type DeviceListCondition struct {
	DeviceNo   string `form:"deviceNo"`
	DeviceType string `form:"deviceType"`
	DeviceInfo string `form:"deviceInfo"`
	Status     string `form:"status"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryDeviceList(ginC *gin.Context) {
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
	orm.Raw("select id from device_info " + conditionStr).Count(&total)

	listSqlStr := "select * from device_info "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "deviceList": rows}})
}
func DeleteDevice(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.DeviceInfo{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyDevice(ginC *gin.Context) {
	var res models.DeviceInfo
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res.Id == 0 {
		db := orm.Save(&models.DeviceInfo{EndoscopeNumber: res.EndoscopeNumber, EndoscopeType: res.EndoscopeType, EndoscopeInfo: res.EndoscopeInfo, Status: res.Status})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	} else {
		db := orm.Save(res)
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type DeviceTypeListCondition struct {
	Name       string `form:"name"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryDeviceTypeList(ginC *gin.Context) {
	var dlc DeviceTypeListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	conditionStr := " where id>0 "
	if len(dlc.Name) > 0 {
		conditionStr += " and name like '%" + dlc.Name + "%' "
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
	orm.Raw("select id from device_type " + conditionStr).Count(&total)

	listSqlStr := "select * from device_type "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "deviceTypeList": rows}})
}
func DeleteDeviceType(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.DeviceType{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyDeviceType(ginC *gin.Context) {
	var res []models.DeviceType
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.DeviceType{Name: res[0].Name})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	} else {
		db := orm.Save(res[0])
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type DoctorListCondition struct {
	Name       string `form:"name"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryDoctorList(ginC *gin.Context) {
	var dlc DoctorListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(dlc.Name) > 0 {
		conditionStr += " and name like '%" + dlc.Name + "%' "
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
	orm.Raw("select id from doctor_info " + conditionStr).Count(&total)

	listSqlStr := "select * from doctor_info "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "doctorList": rows}})
}
func DeleteDoctor(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.DoctorInfo{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyDoctor(ginC *gin.Context) {
	var res []models.DoctorInfo
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.DoctorInfo{Name: res[0].Name})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	} else {
		db := orm.Save(res[0])
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type OperatorListCondition struct {
	Number     string `form:"number"`
	Name       string `form:"name"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryOperatorList(ginC *gin.Context) {
	var dlc OperatorListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(dlc.Number) > 0 {
		conditionStr += " and number like '%" + dlc.Number + "%' "
	}

	if len(dlc.Name) > 0 {
		conditionStr += " and name like '%" + dlc.Name + "%' "
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
	orm.Raw("select id from operator_info " + conditionStr).Count(&total)

	listSqlStr := "select * from operator_info "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "operatorList": rows}})
}
func DeleteOperator(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.OperatorInfo{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyOperator(ginC *gin.Context) {
	var res []models.OperatorInfo
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.OperatorInfo{Number: res[0].Number, Name: res[0].Name})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	} else {
		db := orm.Save(res[0])
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type TimePlanListCondition struct {
	Name              string `form:"name"`
	RinseSolution     string `form:"rinseSolution"`
	DisinfectSolution string `form:"disinfectSolution"`
	Pagination        struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryTimePlanList(ginC *gin.Context) {
	var tpl TimePlanListCondition
	err := ginC.ShouldBind(&tpl)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(tpl.Name) > 0 {
		conditionStr += " and name like '%" + tpl.Name + "%' "
	}

	if len(tpl.RinseSolution) > 0 {
		conditionStr += " and rinse_solution like '%" + tpl.RinseSolution + "%' "
	}

	if len(tpl.DisinfectSolution) > 0 {
		conditionStr += " and disinfect_solution like '%" + tpl.DisinfectSolution + "%' "
	}

	var offset int
	if tpl.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (tpl.Pagination.CurrentPage - 1) * tpl.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from time_plan " + conditionStr).Count(&total)

	listSqlStr := "select * from time_plan "

	// 分页排序
	ordLimitStr := " order by " + tpl.Sort + " " + tpl.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(tpl.Pagination.PageSize), 10)
	log.Print(tpl)
	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": tpl.Pagination.CurrentPage, "pageSize": tpl.Pagination.PageSize, "timeplanList": rows}})
}
func DeleteTimePlan(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.TimePlan{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyTimePlan(ginC *gin.Context) {
	var res []models.TimePlan
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.TimePlan{Name: res[0].Name,
			RinseTime:         res[0].RinseTime,
			RinseSolution:     res[0].RinseSolution,
			DisinfectTime:     res[0].DisinfectTime,
			SteriliseTime:     res[0].SteriliseTime,
			DisinfectSolution: res[0].DisinfectSolution})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	} else {
		db := orm.Save(res[0])
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type ProgramListCondition struct {
	Name       string `form:"name"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func GetProgramSteps(ginC *gin.Context) {
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
	var steps []models.ProgramList
	db := orm.Where("name = ?", res[0]["name"]).Order("id ").Find(&steps)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"steps": &steps}})
}
func GetProgramList(ginC *gin.Context) {
	var tpl ProgramListCondition
	err := ginC.ShouldBind(&tpl)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(tpl.Name) > 0 {
		conditionStr += " and name like '%" + tpl.Name + "%' "
	}

	var offset int
	if tpl.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (tpl.Pagination.CurrentPage - 1) * tpl.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from program " + conditionStr).Count(&total)

	listSqlStr := "select * from program "

	// 分页排序
	ordLimitStr := " order by " + tpl.Sort + " " + tpl.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(tpl.Pagination.PageSize), 10)
	log.Print(tpl)
	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": tpl.Pagination.CurrentPage, "pageSize": tpl.Pagination.PageSize, "programList": rows}})
}
func DeleteProgram(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.Program{}, post[0]["id"])
	orm.Delete(&models.ProgramList{}, "name=?", post[0]["name"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}

type UpdateProgramInfo struct {
	Name  string `json:"Name"`
	TCT   int    `json:"TCT"`
	Steps []struct {
		Index int    `json:"index"`
		Step  string `json:"step"`
		Time  int    `json:"time"`
	} `json:"Steps"`
}

func ModifyProgram(ginC *gin.Context) {
	var upi UpdateProgramInfo
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
	orm.Delete(&models.Program{}, "name=?", upi.Name)
	orm.Delete(&models.ProgramList{}, "name=?", upi.Name)

	orm.Save(&models.Program{Name: upi.Name, TotalCostTime: upi.TCT})
	for _, step := range upi.Steps {
		orm.Save(&models.ProgramList{Name: upi.Name, Step: step.Step, CostTime: step.Time})
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

func GetStepList(ginC *gin.Context) {
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	var step []models.Step
	orm.Find(&step)
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": step})
}

type LiquidDeviceCondition struct {
	DeviceName     string `json:"device_name"`
	DeviceType     string `json:"device_type"`
	DeviceInfo     string `json:"device_info"`
	LiquidType     string `json:"liquid_type"`
	ValidityPeriod int    `json:"validity_period"`
	LimitNumber    int    `json:"limit_number"`
	Pagination     struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryLiquidDeviceList(ginC *gin.Context) {
	var ldc LiquidDeviceCondition
	err := ginC.ShouldBind(&ldc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(ldc.DeviceName) > 0 {
		conditionStr += " and device_name like '%" + ldc.DeviceName + "%' "
	}
	if len(ldc.DeviceType) > 0 {
		conditionStr += " and device_type like '%" + ldc.DeviceType + "%' "
	}
	if len(ldc.DeviceInfo) > 0 {
		conditionStr += " and device_info like '%" + ldc.DeviceInfo + "%' "
	}
	if len(ldc.LiquidType) > 0 {
		conditionStr += " and liquid_type like '%" + ldc.LiquidType + "%' "
	}
	if ldc.ValidityPeriod > 0 {
		conditionStr += " and validity_period <= " + strconv.Itoa(ldc.ValidityPeriod)
	}
	if ldc.LimitNumber > 0 {
		conditionStr += " and limit_number <= " + strconv.Itoa(ldc.LimitNumber)
	}

	var offset int
	if ldc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (ldc.Pagination.CurrentPage - 1) * ldc.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from liquid_device " + conditionStr).Count(&total)

	listSqlStr := "select *, IFNULL(validity_period-timestampdiff(day, (select date from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) ,now()),0) as remaining ,(select usage_count from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) as usagecount from liquid_device as a"

	// 分页排序
	ordLimitStr := " order by " + ldc.Sort + " " + ldc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(ldc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": ldc.Pagination.CurrentPage, "pageSize": ldc.Pagination.PageSize, "LiquidDeviceList": rows}})
}
func DeleteLiquidDeviceList(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据库错误"})
		return
	}
	orm.Delete(&models.LiquidDevice{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyLiquidDeviceList(ginC *gin.Context) {
	var res []models.LiquidDevice
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.LiquidDevice{
			DeviceName:     res[0].DeviceName,
			DeviceType:     res[0].DeviceType,
			DeviceInfo:     res[0].DeviceInfo,
			LiquidType:     res[0].LiquidType,
			ValidityPeriod: res[0].ValidityPeriod,
			LimitNumber:    res[0].LimitNumber,
		})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	} else {
		db := orm.Save(res[0])
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

type LiquidDeviceRecordCondition struct {
	Id         int         `form:"id"`
	Date       []time.Time `form:"date"`
	Pagination struct {
		Total       int  `form:"total"`
		PageSize    int  `form:"pageSize"`
		CurrentPage int  `form:"currentPage"`
		Background  bool `form:"background"`
	} `form:"pagination"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func QueryLiquidDeviceRecordList(ginC *gin.Context) {
	var ldrc LiquidDeviceRecordCondition
	err := ginC.ShouldBind(&ldrc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	ldrc.Date[0] = ldrc.Date[0].In(time.Local)
	ldrc.Date[1] = ldrc.Date[1].In(time.Local)
	var offset int
	if ldrc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (ldrc.Pagination.CurrentPage - 1) * ldrc.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Model(&models.LiquidDetectionRecord{}).Select("id").
		Where("date BETWEEN ? AND ? AND deviceid=?", ldrc.Date[0], ldrc.Date[1], ldrc.Id).Count(&total)

	var rows []models.LiquidDetectionRecord
	orm.Where("date BETWEEN ? AND ? AND deviceid=?", ldrc.Date[0], ldrc.Date[1], ldrc.Id).
		Order(ldrc.Sort + " " + ldrc.Order).Offset(offset).Limit(ldrc.Pagination.PageSize).Find(&rows)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total,
		"currentPage": ldrc.Pagination.CurrentPage, "pageSize": ldrc.Pagination.PageSize, "LiquidDeviceRecordList": rows}})
}

func ModifyLiquidDeviceRecord(ginC *gin.Context) {
	var res models.LiquidDetectionRecord
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}

	orm, _ := models.NewOrm()
	if res.Id == 0 {
		UsageCount := res.UsageCount
		if res.OperatorType == "更换" {
			UsageCount = res.LimitNumber
		}
		db := orm.Save(&models.LiquidDetectionRecord{
			Deviceid:       res.Deviceid,
			LiquidType:     res.LiquidType,
			ValidityPeriod: res.ValidityPeriod,
			LimitNumber:    res.LimitNumber,
			Date:           res.Date,
			Operator:       res.Operator,
			Conclusion:     res.Conclusion,
			OperatorType:   res.OperatorType,
			OperatorInfo:   res.OperatorInfo,
			UsageCount:     UsageCount})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
		ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"msg": "新建成功", "newId": (db.Statement.Model).(*models.LiquidDetectionRecord).Id}})
		return
	} else {
		db := orm.Save(res)
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}

func DeleteLiquidDetection(ginC *gin.Context) {
	var post []map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	orm, err := models.NewOrm()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "数据库错误"})
		return
	}
	orm.Delete(&models.LiquidDetectionRecord{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "删除成功"})
}

func UploadDetectionImg(ginC *gin.Context) {
	Id := ginC.Query("id")
	deviceId := ginC.Query("deviceId")

	f, err := ginC.FormFile("file")
	if err != nil {
		ginC.JSON(200, gin.H{"code": 400, "msg": "上传失败!"})
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			ginC.JSON(200, gin.H{"code": 400, "msg": "上传失败!只允许png,jpg,gif,jpeg文件"})
			return
		}
		path := common.Config.GetString("ecds.detectionimages") + "/" + deviceId + "/" + Id + "/" + time.Now().Format("20060102-150405") + ".jpg"
		if !common.Exist(common.Config.GetString("ecds.detectionimages") + "/" + deviceId) {
			common.TermExec("mkdir", []string{common.Config.GetString("ecds.detectionimages") + "/" + deviceId})
		}
		if !common.Exist(common.Config.GetString("ecds.detectionimages") + "/" + deviceId + "/" + Id) {
			common.TermExec("mkdir", []string{common.Config.GetString("ecds.detectionimages") + "/" + deviceId + "/" + Id})
		}
		log.Print(path)
		ginC.SaveUploadedFile(f, path)
		//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
		ginC.JSON(200, gin.H{"code": 200, "msg": "上传成功!", "result": gin.H{"path": "/images/" +
			deviceId + "/" + Id + "/" + time.Now().Format("20060102-150405") + ".jpg"}})
	}
}

type FileList struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func GetDetectionImagesList(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}

	deviceId := post["deviceId"].(string)
	Id := post["Id"].(string)

	// 获取所有文件
	files, _ := os.ReadDir(common.Config.GetString("ecds.detectionimages") + "/" + deviceId + "/" + Id)
	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	var fileList []FileList
	for _, file := range files {
		if file.IsDir() {
			continue
		} else if file.Name() != ".gitkeep" && file.Name() != ".DS_Store" && file.Name() != ".git" {
			fileList = append(fileList, FileList{Name: file.Name(), Url: "/images/" + deviceId + "/" + Id + "/" + file.Name()})
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": fileList})
}
func DeleteDetectionImages(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "参数错误"})
		return
	}
	filename := post["file"].(string)
	deviceId := post["deviceId"].(string)
	Id := post["Id"].(string)

	common.TermExec("rm", []string{"-rf", common.Config.GetString("ecds.detectionimages") + "/" + deviceId + "/" + Id + "/" + filename})

	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "删除成功！"})
}
func DetectionExport(ginC *gin.Context) {
	var ldrc LiquidDeviceRecordCondition
	err := ginC.ShouldBind(&ldrc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	ldrc.Date[0] = ldrc.Date[0].In(time.Local)
	ldrc.Date[1] = ldrc.Date[1].In(time.Local)

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var device []models.LiquidDevice
	db := orm.Where("id=?", ldrc.Id).Find(&device)

	var rows []models.LiquidDetectionRecord
	db = orm.Where("date BETWEEN ? AND ? AND deviceid=?", ldrc.Date[0], ldrc.Date[1], ldrc.Id).
		Order(ldrc.Sort + " " + ldrc.Order).Find(&rows)

	if db.Error != nil || len(rows) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "没有数据"})
		return
	}
	_ = common.TermExec("rm", []string{"-rf", "./bgStatic/LiquidDetectionExport2*"})
	excelFileName := "./bgStatic/LiquidDetectionExport" + time.Now().Format("2006-01-02") + ".xlsx"
	_ = common.CopyFile("./bgStatic/DetectionExportTemplates.xlsx", excelFileName)
	file, err := excelize.OpenFile(excelFileName)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": "生成文件错误"})
		return
	}
	_ = file.SetCellValue("消毒液检测记录", "C4", device[0].DeviceName)
	_ = file.SetCellValue("消毒液检测记录", "H4", device[0].DeviceType)
	_ = file.SetCellValue("消毒液检测记录", "C6", device[0].DeviceInfo)
	_ = file.SetCellValue("消毒液检测记录", "H6", device[0].LiquidType)

	_ = file.SetCellValue("消毒液检测记录", "C8", time.Now())
	_ = file.SetCellValue("消毒液检测记录", "E8", len(rows))

	rowindex := 10
	for _, val := range rows {
		_ = file.SetCellValue("消毒液检测记录", "A"+strconv.Itoa(rowindex), val.Id)
		_ = file.SetCellValue("消毒液检测记录", "B"+strconv.Itoa(rowindex), val.Deviceid)
		_ = file.SetCellValue("消毒液检测记录", "C"+strconv.Itoa(rowindex), val.LiquidType)
		_ = file.SetCellValue("消毒液检测记录", "D"+strconv.Itoa(rowindex), val.ValidityPeriod)
		_ = file.SetCellValue("消毒液检测记录", "E"+strconv.Itoa(rowindex), val.LimitNumber)
		_ = file.SetCellValue("消毒液检测记录", "F"+strconv.Itoa(rowindex), val.UsageCount)
		_ = file.SetCellValue("消毒液检测记录", "G"+strconv.Itoa(rowindex), val.Date)
		_ = file.SetCellValue("消毒液检测记录", "H"+strconv.Itoa(rowindex), val.Conclusion)
		_ = file.SetCellValue("消毒液检测记录", "I"+strconv.Itoa(rowindex), val.Operator)
		_ = file.SetCellValue("消毒液检测记录", "J"+strconv.Itoa(rowindex), val.OperatorType)
		_ = file.SetCellValue("消毒液检测记录", "K"+strconv.Itoa(rowindex), val.OperatorInfo)

		rowindex++
	}
	_ = file.Save()

	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/LiquidDetectionExport" +
		time.Now().Format("2006-01-02") + ".xlsx"})
}
