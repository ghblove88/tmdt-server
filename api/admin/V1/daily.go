package V1

import (
	"TmdtServer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	if res.ID == 0 {
		db := orm.Save(&models.DeviceInfo{DeviceCode: res.DeviceCode, DeviceSequence: res.DeviceSequence, DeviceName: res.DeviceName, DeviceInfo: res.DeviceInfo})
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
