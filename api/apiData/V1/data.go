package V1

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"EcdsServer/runtime"
	"github.com/gin-gonic/gin"
	"net/http"

	"log"
	"strings"
	"time"
)

func GetOperator(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Operator_Info.GetAll())
}

func GetDoctor(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Doctor_Info.GetAll())
}

func GetDevice(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Device_Info.GetAll())
}

func GetProgram(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Program_Info.GetAll())
}

func GetTimePlan(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_TimePlan_Info.GetAll())
}

// GetLastRecordByEid 返回指定内镜编号 的最后一次洗消记录
func GetLastRecordByEid(ginC *gin.Context) {

	_number := ginC.Param("number")
	//示例中文URL 将 "%AB" 形式的每个 3 字节编码子串转换为 hex-decoded 字节 0xAB
	//_patient, _ = url.QueryUnescape(_patient)f

	// 查询数据 如果 id长度大于4 则是内镜卡号，小于等于4则是内镜身编号
	SqlString := ""
	if len(_number) > 4 {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and endoscope_number='" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	} else {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and Endoscope_info like '%" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	}
	datares := runtime.GetAntRecord(SqlString)

	// 什么也没有找到
	if len(datares) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "count": len(datares), "ant": datares})
}

// GetLastRecordByEidSlic 返回指定内镜编号 的最后一次洗消记录 slic
func GetLastRecordByEidSlic(ginC *gin.Context) {

	_number := ginC.Param("number")

	SqlString := ""
	if len(_number) > 4 {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and endoscope_number='" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	} else {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and Endoscope_info like '%" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	}
	datares := runtime.GetAntRecordSlic(SqlString)

	// 什么也没有找到
	if len(datares) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "count": len(datares), "ant": datares})
}

// GetLastRecordByEidUsed 返回指定内镜编号 的最后一次洗消记录 包含使用过的slic
func GetLastRecordByEidUsed(ginC *gin.Context) {

	_number := ginC.Param("number")

	// 查询数据 如果 id长度大于4 则是内镜卡号，小于等于4则是内镜身编号
	SqlString := ""
	if len(_number) > 4 {
		SqlString = "select * from ant where begin_time <> end_time  and endoscope_number='" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	} else {
		SqlString = "select * from ant where begin_time <> end_time  and Endoscope_info like '%" + _number +
			"'  ORDER BY ID DESC  limit 0,1 "
	}
	datares := runtime.GetAntRecordSlicUsed(SqlString)

	// 什么也没有找到
	if len(datares) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "count": len(datares), "ant": datares})
}

func GetLastRecordByEid_AntListener(ginC *gin.Context) {

	_number := ginC.Param("number")

	antdb, antstepdb := GLRBE(_number)
	//var steps []Steps
	//steps = append(steps, Steps{XMLName: xml.Name{"s", "l"}, Name: "步骤一", CostTime: 12, Wm: "洗净机1"})
	//steps = append(steps, Steps{XMLName: xml.Name{"s", "l"}, Name: "步骤二", CostTime: 12, Wm: "洗净机1"})
	//steps = append(steps, Steps{XMLName: xml.Name{"s", "l"}, Name: "步骤三", CostTime: 12, Wm: "洗净机1"})

	if antdb.Id == 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false})
		return
	} else {

		db := gin.H{"Id": antdb.Id, "Number": antdb.Number, "EndoscopeNumber": antdb.EndoscopeNumber,
			"EndoscopeType": antdb.EndoscopeType, "Operator": antdb.Operator, "PatientName": antdb.PatientName,
			"DocName": antdb.DocName, "Diseases": antdb.Diseases,
			"BeginTime":     antdb.BeginTime.Format("2006-01-02 15:04:05"),
			"EndTime":       antdb.EndTime.Format("2006-01-02 15:04:05"),
			"TotalCostTime": antdb.TotalCostTime, "EndoscopeInfo": antdb.EndoscopeInfo,
			"Endoscope_Model": strings.Split(antdb.EndoscopeInfo, "/")[0]}

		ginC.JSON(http.StatusOK, gin.H{"success": true, "id": antdb.Id, "ant": db, "step": antstepdb})
	}
}

func GLRBE(msg string) (antdb models.Ant, antstepdb []models.AntStep) {

	engine, err := models.NewOrm()

	if err != nil {
		println(err.Error())
		return antdb, antstepdb
	}

	var ant models.Ant
	if len(msg) > 4 {
		engine.Where("patient_name = '' and begin_time <> end_time and to_days(begin_time) = to_days(now()) and endoscope_number = ?", msg).Order("id desc").First(&ant)
	} else {
		engine.Where("patient_name = '' and begin_time <> end_time and to_days(begin_time) = to_days(now()) and Endoscope_info LIKE ?", "%"+msg+"%").Order("id desc").First(&ant)
	}

	if ant.Id > 0 {
		if antdb.Id == 0 {
			antdb = ant
			engine.Where("number = ?", antdb.Number).Find(&antstepdb)
		} else {
			if antdb.EndTime.UnixNano() < ant.EndTime.UnixNano() {
				antdb = ant
				antstepdb = []models.AntStep{}
				engine.Where("number = ?", antdb.Number).Find(&antstepdb)
			}
		}
	}

	if err != nil {
		return models.Ant{}, nil
	}

	return antdb, antstepdb
}

// GetRecordByNo 返回指定洗消记录编号的记录
func GetRecordByNo(ginC *gin.Context) {
	_number := ginC.Param("number")

	SqlString := "select * from ant where begin_time <> end_time and number='" + _number +
		"'  ORDER BY ID DESC  limit 0,1 "
	ginC.JSON(http.StatusOK, gin.H{"success": true, "ant": runtime.GetAntRecord(SqlString)})
}

// GetRecordsById 返回指定ID 后的所有记录
func GetRecordsById(ginC *gin.Context) {
	_number := ginC.Param("number")
	SqlString := "select * from ant where begin_time <> end_time and id>" + _number +
		" ORDER BY ID"
	ginC.JSON(http.StatusOK, gin.H{"success": true, "ant": runtime.GetAntRecord(SqlString)})
}

// SwipeCard 模拟刷卡
func SwipeCard(ginC *gin.Context) {
	ip := strings.Split(ginC.RemoteIP(), ":")[0]
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	if post["ip"] != nil {
		ip = post["ip"].(string)
	}
	runtime.G_Socket_Reader.Readid_queue.Push_ReadID_Queue(
		runtime.STRUCT_READID_MSG{ip, post["number"].(string)})

	log.Println("URL 刷卡操作:", runtime.G_Reader_Info.Find(ip).Name, ip, post["number"].(string))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "刷卡成功"})
}

// WriteBack1 回写病人信息
func WriteBack1(ginC *gin.Context) {
	var dat models.XianWriteback
	if err := ginC.BindJSON(dat); err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据解析失败"})
		return
	}

	if len(dat.Rfidno) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "无效卡号"})
		return
	}

	// 解除绑定
	if dat.Bindmirrostate == "0" {
		if len(dat.Number) <= 0 {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解绑失败！无效洗消编号。"})
			return
		}

		sqlstr := "SELECT number FROM xian_writeback WHERE rfidNo='" + dat.Rfidno + "' and patientId='" + dat.Patientid + "' and checkReportId='" + dat.Checkreportid + "'"
		var wlog []map[string]interface{}
		o, _ := models.NewOrm()
		o.Raw(sqlstr).Scan(&wlog)
		if len(wlog) > 0 {
			sqlstr := "update ant set patient_name='' where number='" + wlog[0]["number"].(string) + "'"
			db := o.Exec(sqlstr)
			if db.Error != nil {
				log.Println(db.Error)
				ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解绑失败！"})
				return
			}
			//记录回写过来的完整数据
			dat.Number = wlog[0]["number"].(string)
			db = o.Select("number").Create(&wlog)
			if db.Error != nil {
				log.Println(db.Error)
			}
			ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "解绑成功"})
			return
		}

		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解绑失败，未找到绑定记录！"})
		return
	}

	// 查询数据 如果 id长度大于4 则是内镜卡号，小于等于4则是内镜身编号
	SqlString := ""
	if len(dat.Rfidno) > 4 {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and endoscope_number='" + dat.Rfidno +
			"'  ORDER BY ID DESC  limit 0,1 "
	} else {
		SqlString = "select * from ant where begin_time <> end_time and patient_name = '' and Endoscope_info like '%" + dat.Rfidno +
			"'  ORDER BY ID DESC  limit 0,1 "
	}
	datares := runtime.GetAntRecordSlic(SqlString)

	if len(datares) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "未找到绑定记录！"})
		return
	}

	if len(dat.Namepatient) > 0 {
		sqlstr := "update ant set patient_name='" + dat.Namepatient + "' where number='" + datares[0].A_number + "' and patient_name=''"
		o, _ := models.NewOrm()
		db := o.Exec(sqlstr)
		if db.Error != nil {
			log.Println(db.Error)
		}
		log.Println("更新洗消记录", datares[0].A_number, "病人信息为:"+dat.Namepatient)

		//记录回写过来的完整数据
		dat.Number = datares[0].A_number
		db = o.Select("number").Create(&dat)
		if db.Error != nil {
			log.Println(db.Error)
		}
	}

	ginC.JSON(http.StatusOK, gin.H{"code": 1, "success": true, "msg": "成功",
		"data": gin.H{"scopeNo": datares[0].A_endoscope_info, "scopeType": datares[0].A_endoscope_type,
			"scopeInnerNo": common.GetDeviceVoice(dat.Rfidno, datares[0].A_endoscope_info), "rfidNo": dat.Rfidno}})

}

// WriteBack2 回写病人信息
func WriteBack2(ginC *gin.Context) {
	var dat models.DongguanWriteback
	if err := ginC.BindJSON(&dat); err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据解析失败！"})
		return
	}

	if len(dat.Rfidno) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "无效卡号！"})
		return
	}

	// 解除绑定
	if dat.Bindmirrostate == "0" {

		if len(dat.Number) <= 0 {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解绑失败，无效洗消编号！"})
			return
		}

		o, _ := models.NewOrm()
		sqlstr := "update ant set patient_name='' where number='" + dat.Number + "'"
		db := o.Exec(sqlstr)
		if db.Error != nil {
			log.Println(db.Error)
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解绑失败！"})
			return
		}
		//记录回写过来的完整数据
		o.Where("number=?", dat.Number).Delete(&models.DongguanWriteback{})

		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "解绑成功！"})
		return
	}

	SqlString := "select * from ant where patient_name='' and number='" + dat.Number + "'  ORDER BY ID DESC  limit 0,1 "

	datares := runtime.GetAntRecordSlic(SqlString)

	if len(datares) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "记录已经被绑定，或未找到有效记录"})
		return
	}

	if len(dat.Namepatient) > 0 {
		sqlstr := "update ant set patient_name='" + dat.Namepatient + "' where number='" + datares[0].A_number + "' and patient_name=''"
		o, _ := models.NewOrm()
		db := o.Exec(sqlstr)
		if db.Error != nil {
			log.Println(db.Error)
		}
		log.Println("更新洗消记录", datares[0].A_number, "病人信息为:"+dat.Namepatient)

		//记录回写过来的完整数据
		dat.Number = datares[0].A_number
		db = o.Create(&dat)
		if db.Error != nil {
			log.Println(db.Error)
		}
	}

	ginC.JSON(http.StatusOK, gin.H{"code": 1, "success": true, "msg": "成功",
		"data": gin.H{"endoscopeInfo": datares[0].A_endoscope_info, "endoscopeType": datares[0].A_endoscope_type,
			"InnerNo": common.GetDeviceVoice(dat.Rfidno, datares[0].A_endoscope_info), "rfidNo": dat.Rfidno}})
}

func WriteBack2_AntListener(ginC *gin.Context) {
	var dat models.DongguanWriteback
	if err := ginC.BindJSON(&dat); err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据解析失败！"})
		return
	}

	if len(dat.Rfidno) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "无效卡号！"})
		return
	}

	// 解除绑定
	if dat.Bindmirrostate == "0" {
		engine, _ := models.NewOrm()
		var logdb models.DongguanWriteback
		engine.Where("number = ?", dat.Number).First(logdb)

		if logdb.Id > 0 {
			engine.Exec("update ant set patient_name='' where number='" + logdb.Number + "'")
			engine.Where("number = ?", dat.Number).Delete(&models.DongguanWriteback{})

			ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "解除绑定成功！"})
			return
		}

		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "未找到可解绑的记录！"})
		return
	}

	if len(dat.Namepatient) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "病人姓名不能为空！"})
		return
	}

	antdb, _ := GLRBE(dat.Rfidno)
	if antdb.Id <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "未找到有效记录！"})
		return
	}

	sqlstr := "update ant set patient_name='" + dat.Namepatient + "' where number='" + antdb.Number + "' and patient_name=''"
	engine, _ := models.NewOrm()
	engine.Exec(sqlstr)
	dat.Number = antdb.Number
	dat.BeginTime = time.Now()
	engine.Create(dat)

	ginC.JSON(http.StatusOK, gin.H{"code": 1, "success": true, "msg": "成功",
		"data": gin.H{"scopeNo": antdb.EndoscopeInfo, "scopeType": antdb.EndoscopeType,
			"scopeInnerNo": common.GetDeviceVoice(dat.Rfidno, antdb.EndoscopeInfo), "rfidNo": antdb.EndoscopeNumber}})
	return
}

// EndoscopyWriteBack 使用后记录模式，先记录工作站检查病人信息，再清洗镜子是绑定清洗信息
func EndoscopyWriteBack(ginC *gin.Context) {
	var dat models.EndoscopyWriteback
	if err := ginC.BindJSON(&dat); err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据解析失败！"})
		return
	}

	if len(dat.Rfidno) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "无效卡号！"})
		return
	}

	//记录回写过来的完整数据
	dat.Number = ""
	orm, _ := models.NewOrm()
	db := orm.Create(&dat)
	if db.Error != nil {
		log.Println(db.Error)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "记录失败，唯一标识重复！"})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "信息记录成功！"})
}
