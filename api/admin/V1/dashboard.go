package V1

import (
	"EcdsServer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetStatistics1(ginC *gin.Context) {
	var timeLimit map[string][]time.Time
	err := ginC.ShouldBind(&timeLimit)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	timeLimit["date"][0] = timeLimit["date"][0].In(time.Local)
	timeLimit["date"][1] = timeLimit["date"][1].In(time.Local)

	var antTotal int64
	orm, _ := models.NewOrm()
	orm.Model(&models.Ant{}).Where("begin_time between ? AND ?", timeLimit["date"][0], timeLimit["date"][1]).Count(&antTotal)

	var row []map[string]interface{}
	orm.Raw("SELECT sum(total_cost_time) as totalTime FROM ant WHERE begin_time between ? AND ?",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&row)

	var totalTime int
	if len(row) > 0 && row[0]["totalTime"] != nil {
		totalTime, _ = strconv.Atoi(row[0]["totalTime"].(string))
		totalTime = totalTime / 60 / 60
	}

	row = []map[string]interface{}{}
	orm.Raw("SELECT diseases,count(id) as count FROM ant WHERE begin_time between ? AND ? GROUP BY diseases ORDER BY diseases",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&row)

	var countDiseases int64
	var countNonDiseases int64
	if len(row) > 1 {
		countNonDiseases, _ = row[1]["count"].(int64)
	}
	if len(row) > 2 {
		countDiseases, _ = row[2]["count"].(int64)
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功",
		"data": gin.H{"antTotal": antTotal,
			"totalTime":        totalTime,
			"countDiseases":    countDiseases,
			"countNonDiseases": countNonDiseases,
		}})
}

func GetStatistics2(ginC *gin.Context) {
	var timeLimit map[string][]time.Time
	err := ginC.ShouldBind(&timeLimit)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	timeLimit["date"][0] = timeLimit["date"][0].In(time.Local)
	timeLimit["date"][1] = timeLimit["date"][1].In(time.Local)

	var row []map[string]interface{}
	orm, _ := models.NewOrm()
	orm.Raw("SELECT step,count(step) as count,SUM(cost_time) as total_time FROM ant_step LEFT JOIN ant ON ant_step.number=ant.number WHERE begin_time between  ? AND ? GROUP BY step ORDER BY step",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&row)

	if len(row) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"dimensions": []struct {
			product   string
			count     int
			totalTime int
		}{}, "source": []string{}}})
		return
	}

	var dimensions = []string{"步骤", "次数", "时长(小时)"}
	var source []struct {
		Product   string `json:"步骤"`
		Count     int    `json:"次数"`
		TotalTime int    `json:"时长(小时)"`
	}
	for _, v := range row {
		totalTime, _ := strconv.Atoi(v["total_time"].(string))
		source = append(source, struct {
			Product   string `json:"步骤"`
			Count     int    `json:"次数"`
			TotalTime int    `json:"时长(小时)"`
		}{Product: v["step"].(string), Count: int(v["count"].(int64)) / 10, TotalTime: totalTime / 60 / 60})
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"dimensions": dimensions, "source": source}})
}

func GetStatistics3(ginC *gin.Context) {
	var timeLimit map[string][]time.Time
	err := ginC.ShouldBind(&timeLimit)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	timeLimit["date"][0] = timeLimit["date"][0].In(time.Local)
	timeLimit["date"][1] = timeLimit["date"][1].In(time.Local)

	var row []map[string]interface{}
	orm, _ := models.NewOrm()
	orm.Raw("SELECT washing_machine as wm, COUNT(DISTINCT(ant.number)) as count,SUM(ant_step.cost_time) as total_time FROM ant LEFT JOIN ant_step ON ant.number = ant_step.number  WHERE begin_time between  ? AND ? GROUP BY washing_machine",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&row)

	if len(row) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"dimensions": []struct {
			product   string
			count     int
			totalTime int
		}{}, "source": []string{}}})
		return
	}

	var dimensions = []string{"洗净机名称", "次数", "时长(小时)"}
	var source []struct {
		Product   string `json:"洗净机名称"`
		Count     int    `json:"次数"`
		TotalTime int    `json:"时长(小时)"`
	}
	for _, v := range row {
		if v["wm"] == nil || v["wm"].(string) == "" {
			continue
		}
		totalTime, _ := strconv.Atoi(v["total_time"].(string))
		source = append(source, struct {
			Product   string `json:"洗净机名称"`
			Count     int    `json:"次数"`
			TotalTime int    `json:"时长(小时)"`
		}{Product: v["wm"].(string), Count: int(v["count"].(int64)) / 10, TotalTime: totalTime / 60 / 60})
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"dimensions": dimensions, "source": source}})
}

func GetStatistics4(ginC *gin.Context) {
	var timeLimit map[string][]time.Time
	err := ginC.ShouldBind(&timeLimit)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	timeLimit["date"][0] = timeLimit["date"][0].In(time.Local)
	timeLimit["date"][1] = timeLimit["date"][1].In(time.Local)

	var doctorRow []map[string]interface{}
	orm, _ := models.NewOrm()
	orm.Raw("SELECT doc_name ,COUNT(doc_name) as count,SUM(total_cost_time) as cost_time,COUNT(if(ant.diseases=2,1,NULL)) as diseases1 FROM ant WHERE begin_time between  ? AND ? GROUP BY doc_name",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&doctorRow)

	var operatorRow []map[string]interface{}
	orm.Raw("SELECT operator ,COUNT(operator) as count,SUM(total_cost_time) as cost_time,COUNT(if(ant.diseases=2,1,NULL)) as diseases1 FROM ant WHERE begin_time between  ? AND ? GROUP BY operator",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&operatorRow)

	var endoscopeRow []map[string]interface{}
	orm.Raw("SELECT endoscope_number as number,COUNT(endoscope_number) as count,SUM(total_cost_time) as cost_time,COUNT(if(ant.diseases=2,1,NULL)) as diseases1 FROM ant WHERE begin_time between  ? AND ? GROUP BY endoscope_number",
		timeLimit["date"][0], timeLimit["date"][1]).Find(&endoscopeRow)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"doctor": doctorRow, "operator": operatorRow, "endoscope": endoscopeRow}})
}

func GetModule1(ginC *gin.Context) {
	var antCount int64
	var doctorCount int64
	var operatorCount int64
	var deviceCount int64
	orm, _ := models.NewOrm()
	orm.Model(&models.Ant{}).Count(&antCount)
	orm.Model(&models.DoctorInfo{}).Count(&doctorCount)
	orm.Model(&models.OperatorInfo{}).Count(&operatorCount)
	orm.Model(&models.DeviceInfo{}).Count(&deviceCount)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功",
		"data": gin.H{"antCount": antCount,
			"doctorCount":   doctorCount,
			"operatorCount": operatorCount,
			"deviceCount":   deviceCount,
		}})
}

func GetModule2(ginC *gin.Context) {
	var QueryType map[string]string
	err := ginC.ShouldBind(&QueryType)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	if QueryType["data"] == "近七日" {
		orm, _ := models.NewOrm()
		var category [7]string
		var data [7]int64
		for i := 7; i > 0; i-- {
			var begin = time.Now().AddDate(0, 0, -(i - 1))
			var end = time.Now().AddDate(0, 0, -(i - 1))
			var count int64
			orm.Model(&models.Ant{}).Where("begin_time between ? AND ?", begin.Format("2006-01-02 ")+"00:00:01", end.Format("2006-01-02 ")+"23:59:59").Count(&count)
			category[7-i] = begin.Format("01-02")
			data[7-i] = count
		}
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"category": category, "data": data}})
		return
	}
	if QueryType["data"] == "近三十天" {
		orm, _ := models.NewOrm()
		var category [30]string
		var data [30]int64
		for i := 30; i > 0; i-- {
			time.Now().Day()
			var begin = time.Now().AddDate(0, 0, -(i - 1))
			var end = time.Now().AddDate(0, 0, -(i - 1))
			var count int64
			orm.Model(&models.Ant{}).Where("begin_time between ? AND ?", begin.Format("2006-01-02 ")+"00:00:01", end.Format("2006-01-02 ")+"23:59:59").Count(&count)
			category[30-i] = begin.Format("01-02")
			data[30-i] = count
		}
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"category": category, "data": data}})
		return
	}
	if QueryType["data"] == "近一季度" {
		orm, _ := models.NewOrm()
		var category [12]string
		var data [12]int64
		for i := 12; i > 0; i-- {
			time.Now().Day()
			var begin = time.Now().AddDate(0, 0, -(i*7)-(i-1))
			var end = time.Now().AddDate(0, 0, -(i-1)*7-(i-1))
			var count int64
			orm.Model(&models.Ant{}).Where("begin_time between ? AND ?", begin.Format("2006-01-02 ")+"00:00:01", end.Format("2006-01-02 ")+"23:59:59").Count(&count)
			category[12-i] = begin.Format("01/02") + end.Format("-02")
			data[12-i] = count
		}
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"category": category, "data": data}})
		return
	}
	if QueryType["data"] == "近一年" {
		orm, _ := models.NewOrm()
		var category [12]string
		var data [12]int64
		for i := 12; i > 0; i-- {
			time.Now().Day()
			var begin = time.Now().AddDate(0, -i+1, -time.Now().Day()+1)
			var end = time.Now().AddDate(0, -i+2, -time.Now().Day())
			var count int64
			orm.Model(&models.Ant{}).Where("begin_time between ? AND ?", begin.Format("2006-01-02 ")+"00:00:01", end.Format("2006-01-02 ")+"23:59:59").Count(&count)
			category[12-i] = begin.Format("2006/01")
			data[12-i] = count
		}
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"category": category, "data": data}})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
	return
}

func GetModule3(ginC *gin.Context) {
	var QueryType map[string]string
	err := ginC.ShouldBind(&QueryType)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var deviceName []string
	var Remaining []int64
	var usageCount []int64

	sql := "select *, IFNULL(validity_period-timestampdiff(day, (select date from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) ,now()),0) as remaining ,(select usage_count from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) as usagecount from liquid_device as a WHERE device_type<>'暂停使用'"
	var rows []map[string]interface{}
	orm, err := models.NewOrm()

	db := orm.Raw(sql).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	for _, v := range rows {
		deviceName = append(deviceName, v["device_name"].(string))
		Remaining = append(Remaining, v["remaining"].(int64))
		usageCount = append(usageCount, v["usagecount"].(int64))
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": gin.H{"Name": deviceName, "Remaining": Remaining, "usageCount": usageCount}})
	return
}

func GetModule4(ginC *gin.Context) {

	var row []map[string]interface{}
	orm, _ := models.NewOrm()
	db := orm.Raw("SELECT endoscope_type as name,COUNT(endoscope_type) as value FROM device_info GROUP BY endoscope_type").Find(&row)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取成功", "data": row})

}
