package V1

import (
	"TmdtServer/models"
	"TmdtServer/runtime"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DeviceListCondition struct {
	DeviceCode     string `json:"device_code"`
	DeviceSequence string `json:"device_sequence"`
	DeviceName     string `json:"device_name"`
	DeviceInfo     string `json:"device_info"`
	Pagination     struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryDeviceList(ginC *gin.Context) {
	var dlc DeviceListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(dlc.DeviceCode) > 0 {
		conditionStr += " and device_code like '%" + dlc.DeviceCode + "%' "
	}
	if len(dlc.DeviceSequence) > 0 {
		conditionStr += " and device_sequence like '%" + dlc.DeviceSequence + "%' "
	}
	if len(dlc.DeviceName) > 0 {
		conditionStr += " and device_name like '%" + dlc.DeviceName + "%' "
	}
	if len(dlc.DeviceInfo) > 0 {
		conditionStr += " and device_info like '%" + dlc.DeviceInfo + "%' "
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
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res.ID == 0 {
		db := orm.Save(&models.DeviceInfo{DeviceCode: res.DeviceCode, DeviceSequence: res.DeviceSequence, DeviceName: res.DeviceName, DeviceInfo: res.DeviceInfo})
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	} else {
		db := orm.Save(res)
		if db.Error != nil {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
			return
		}
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "修改成功"})
}
func QueryDeviceCode(ginC *gin.Context) {

	var rows []struct{ Code string }
	for _, device := range runtime.G_SocketServer.DataMap {
		rows = append(rows, struct{ Code string }{Code: strconv.Itoa(int(device.DeviceID))})
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"deviceCodeList": rows}})
}

type BedsListCondition struct {
	BedCode    string `json:"bed_code"`
	Pagination struct {
		Total       int  `json:"total"`
		PageSize    int  `json:"pageSize"`
		CurrentPage int  `json:"currentPage"`
		Background  bool `json:"background"`
	} `json:"pagination"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func QueryBedsList(ginC *gin.Context) {
	var dlc BedsListCondition
	err := ginC.ShouldBind(&dlc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(dlc.BedCode) > 0 {
		conditionStr += " and bed_code like '%" + dlc.BedCode + "%' "
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
	orm.Raw("select id from bed " + conditionStr).Count(&total)

	listSqlStr := "select * from bed "

	// 分页排序
	ordLimitStr := " order by " + dlc.Sort + " " + dlc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(dlc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"total": total, "currentPage": dlc.Pagination.CurrentPage, "pageSize": dlc.Pagination.PageSize, "bedList": rows}})
}
func DeleteBed(ginC *gin.Context) {
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
	orm.Delete(&models.Bed{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}
func ModifyBed(ginC *gin.Context) {
	var res []models.Bed
	err := ginC.ShouldBind(&res)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].ID == 0 {
		db := orm.Save(&models.Bed{BedCode: res[0].BedCode, BedInfo: res[0].BedInfo})
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
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "修改成功"})
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
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	if res[0].Id == 0 {
		db := orm.Save(&models.OperatorInfo{Number: res[0].Number, Name: res[0].Name})
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
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "修改成功"})
}
