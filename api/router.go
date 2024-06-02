package api

import (
	adminDataV1 "EcdsServer/api/admin/V1"
	apiDataV1 "EcdsServer/api/apiData/V1"
	monitorV1 "EcdsServer/api/monitor/V1"
	"EcdsServer/common"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Router() {
	gin.SetMode(gin.ReleaseMode)
	logfile, _ := os.Create("./http.log")
	gin.DefaultWriter = io.MultiWriter(logfile)

	r := gin.Default()
	//r.Use(common.GinLogger())

	err := r.SetTrustedProxies(nil)
	if err != nil {
		return
	}

	// 中间件，只有在中间件之后注册的路由才会走中间件
	r.Use(cors())

	// 静态文件
	r.Static("/bgStatic", "./bgStatic")
	r.Static("/images", common.Config.GetString("ecds.detectionimages"))

	// 创建基于cookie的存储引擎， 参数是用于加密的密钥
	store := cookie.NewStore([]byte("5e59a3e7-2571-4f9b-81ca-1a7682af3aad"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))

	apiDataV1Group := r.Group("/data/v1")
	{
		apiDataV1Group.POST("/SwipeCard", func(c *gin.Context) { apiDataV1.SwipeCard(c) })

		apiDataV1Group.GET("/operator", func(c *gin.Context) { apiDataV1.GetOperator(c) })
		apiDataV1Group.GET("/doctor", func(c *gin.Context) { apiDataV1.GetDoctor(c) })
		apiDataV1Group.GET("/device", func(c *gin.Context) { apiDataV1.GetDevice(c) })
		apiDataV1Group.GET("/program", func(c *gin.Context) { apiDataV1.GetProgram(c) })
		apiDataV1Group.GET("/timePlan", func(c *gin.Context) { apiDataV1.GetTimePlan(c) })
		apiDataV1Group.GET("/lastEndoscopeDisinfectionRecord/:number", func(c *gin.Context) { apiDataV1.GetLastRecordByEid(c) })
		apiDataV1Group.GET("/lastEndoscopeDisinfectionRecord1/:number", func(c *gin.Context) { apiDataV1.GetLastRecordByEidSlic(c) })
		apiDataV1Group.GET("/lastEndoscopeDisinfectionRecord2/:number", func(c *gin.Context) { apiDataV1.GetLastRecordByEid_AntListener(c) })
		apiDataV1Group.GET("/lastEndoscopeDisinfectionRecord3/:number", func(c *gin.Context) { apiDataV1.GetLastRecordByEidUsed(c) })
		apiDataV1Group.GET("/recordByBo/:number", func(c *gin.Context) { apiDataV1.GetRecordByNo(c) })
		apiDataV1Group.GET("/recordsById/:number", func(c *gin.Context) { apiDataV1.GetRecordsById(c) })

		apiDataV1Group.POST("/writeBack1", func(c *gin.Context) { apiDataV1.WriteBack1(c) })
		apiDataV1Group.POST("/writeBack2", func(c *gin.Context) { apiDataV1.WriteBack2(c) })
		apiDataV1Group.POST("/writeBack2a", func(c *gin.Context) { apiDataV1.WriteBack2_AntListener(c) })
		apiDataV1Group.POST("/endoscopyWriteBack", func(c *gin.Context) { apiDataV1.EndoscopyWriteBack(c) })

	}

	apiAdminV1Group := r.Group("/admin/v1")
	{
		// 服务状态检测提示
		r.GET("/", func(c *gin.Context) { adminDataV1.Welcome(c) })

		apiAdminV1Group.POST("/login", func(c *gin.Context) { adminDataV1.Login(c) })
		apiAdminV1Group.GET("/get-async-routes", func(c *gin.Context) { adminDataV1.GetAsyncRoutes(c) })
		apiAdminV1Group.GET("/deviceType/list", func(c *gin.Context) { adminDataV1.GetDeviceTypeList(c) })
		apiAdminV1Group.GET("/operator/list", func(c *gin.Context) { adminDataV1.GetOperatorList(c) })
		apiAdminV1Group.GET("/doctor/list", func(c *gin.Context) { adminDataV1.GetDoctorList(c) })
		apiAdminV1Group.GET("/ecm/list", func(c *gin.Context) { adminDataV1.GetEcmList(c) })

		user := apiAdminV1Group.Group("/ant")
		{
			user.POST("/antList", func(c *gin.Context) { adminDataV1.QueryAntList(c) })
			user.POST("/antStepList", func(c *gin.Context) { adminDataV1.QueryAntSteps(c) })
			user.POST("/antUpdate", func(c *gin.Context) { adminDataV1.AntModify(c) })
			user.POST("/AntBatchExport", func(c *gin.Context) { adminDataV1.AntBatchExport(c) })
		}

		daily := apiAdminV1Group.Group("/daily")
		{
			daily.POST("/queryDevice", func(c *gin.Context) { adminDataV1.QueryDeviceList(c) })
			daily.POST("/deleteDevice", func(c *gin.Context) { adminDataV1.DeleteDevice(c) })
			daily.POST("/updateDevice", func(c *gin.Context) { adminDataV1.ModifyDevice(c) })

			daily.POST("/queryDeviceType", func(c *gin.Context) { adminDataV1.QueryDeviceTypeList(c) })
			daily.POST("/deleteDeviceType", func(c *gin.Context) { adminDataV1.DeleteDeviceType(c) })
			daily.POST("/updateDeviceType", func(c *gin.Context) { adminDataV1.ModifyDeviceType(c) })

			daily.POST("/queryDoctor", func(c *gin.Context) { adminDataV1.QueryDoctorList(c) })
			daily.POST("/deleteDoctor", func(c *gin.Context) { adminDataV1.DeleteDoctor(c) })
			daily.POST("/updateDoctor", func(c *gin.Context) { adminDataV1.ModifyDoctor(c) })

			daily.POST("/queryOperator", func(c *gin.Context) { adminDataV1.QueryOperatorList(c) })
			daily.POST("/deleteOperator", func(c *gin.Context) { adminDataV1.DeleteOperator(c) })
			daily.POST("/updateOperator", func(c *gin.Context) { adminDataV1.ModifyOperator(c) })

			daily.POST("/queryTimePlan", func(c *gin.Context) { adminDataV1.QueryTimePlanList(c) })
			daily.POST("/deleteTimePlan", func(c *gin.Context) { adminDataV1.DeleteTimePlan(c) })
			daily.POST("/updateTmePlan", func(c *gin.Context) { adminDataV1.ModifyTimePlan(c) })

			daily.POST("/getProgramSteps", func(c *gin.Context) { adminDataV1.GetProgramSteps(c) })
			daily.POST("/queryProgram", func(c *gin.Context) { adminDataV1.GetProgramList(c) })
			daily.POST("/deleteProgram", func(c *gin.Context) { adminDataV1.DeleteProgram(c) })
			daily.POST("/updateProgram", func(c *gin.Context) { adminDataV1.ModifyProgram(c) })
			daily.GET("/getStepList", func(c *gin.Context) { adminDataV1.GetStepList(c) })

			daily.POST("/queryLiquidDeviceList", func(c *gin.Context) { adminDataV1.QueryLiquidDeviceList(c) })
			daily.POST("/deleteLiquidDevice", func(c *gin.Context) { adminDataV1.DeleteLiquidDeviceList(c) })
			daily.POST("/saveLiquidDevice", func(c *gin.Context) { adminDataV1.ModifyLiquidDeviceList(c) })
			daily.POST("/queryLiquidDeviceRecordList", func(c *gin.Context) { adminDataV1.QueryLiquidDeviceRecordList(c) })
			daily.POST("/saveLiquidDeviceRecord", func(c *gin.Context) { adminDataV1.ModifyLiquidDeviceRecord(c) })
			daily.POST("/deleteLiquidDetection", func(c *gin.Context) { adminDataV1.DeleteLiquidDetection(c) })
			daily.POST("/uploadDetectionImg", func(c *gin.Context) { adminDataV1.UploadDetectionImg(c) })
			daily.POST("/getDetectionImagesList", func(c *gin.Context) { adminDataV1.GetDetectionImagesList(c) })
			daily.POST("/deleteDetectionImages", func(c *gin.Context) { adminDataV1.DeleteDetectionImages(c) })
			daily.POST("/detectionExport", func(c *gin.Context) { adminDataV1.DetectionExport(c) })
		}

		endoscope := apiAdminV1Group.Group("/endoscope")
		{
			endoscope.POST("/queryLeakDetectionList", func(c *gin.Context) { adminDataV1.QueryLeakDetectionList(c) })
			endoscope.POST("/savePartsModify", func(c *gin.Context) { adminDataV1.ModifyPartsModify(c) })
			endoscope.POST("/leakDetectionExport", func(c *gin.Context) { adminDataV1.LeakDetectionExport(c) })
			endoscope.POST("/queryStorageList", func(c *gin.Context) { adminDataV1.GetStorageList(c) })
			endoscope.POST("/queryStorageRecordList", func(c *gin.Context) { adminDataV1.GetStorageRecordList(c) })
			endoscope.POST("/storageExport", func(c *gin.Context) { adminDataV1.StorageExport(c) })
			endoscope.POST("/repairRecordExport", func(c *gin.Context) { adminDataV1.RepairRecordExport(c) })
			endoscope.POST("/sendRepairModify", func(c *gin.Context) { adminDataV1.SendRepairModify(c) })
			endoscope.POST("/completedRepairModify", func(c *gin.Context) { adminDataV1.CompletedRepairModify(c) })
			endoscope.POST("/queryRepairDetailsList", func(c *gin.Context) { adminDataV1.QueryRepairDetailsList(c) })
			endoscope.POST("/queryRepairRecordList", func(c *gin.Context) { adminDataV1.QueryRepairRecordList(c) })
			endoscope.POST("/repairManagementExport", func(c *gin.Context) { adminDataV1.RepairManagementExport(c) })
		}

		system := apiAdminV1Group.Group("/system")
		{
			system.POST("/backupData", func(c *gin.Context) { adminDataV1.BackupData(c) })
			system.POST("/backupFileUpload", func(c *gin.Context) { adminDataV1.BackupFileUpload(c) })
			system.POST("/dataRecovery", func(c *gin.Context) { adminDataV1.DataRecovery(c) })

			system.POST("/getUserList", func(c *gin.Context) { adminDataV1.QueryUserList(c) })
			system.POST("/UserDelete", func(c *gin.Context) { adminDataV1.DeleteUser(c) })
			system.POST("/UserModify", func(c *gin.Context) { adminDataV1.UserModify(c) })
			system.POST("/resetPassword", func(c *gin.Context) { adminDataV1.ResetPassword(c) })
			system.POST("/changePassword", func(c *gin.Context) { adminDataV1.ChangePassword(c) })
			system.POST("/restart", func(c *gin.Context) { adminDataV1.Restart(c) })
			system.GET("/getSystemInfo", func(c *gin.Context) { adminDataV1.GetSystemInfo(c) })
			system.GET("/getConfigurationFile", func(c *gin.Context) { adminDataV1.GetConfigurationFile(c) })
			system.POST("/saveConfigurationFile", func(c *gin.Context) { adminDataV1.SaveConfigurationFile(c) })
		}

		dashboard := apiAdminV1Group.Group("/dashboard")
		{
			dashboard.POST("/getStatistics1", func(c *gin.Context) { adminDataV1.GetStatistics1(c) })
			dashboard.POST("/getStatistics2", func(c *gin.Context) { adminDataV1.GetStatistics2(c) })
			dashboard.POST("/getStatistics3", func(c *gin.Context) { adminDataV1.GetStatistics3(c) })
			dashboard.POST("/getStatistics4", func(c *gin.Context) { adminDataV1.GetStatistics4(c) })

			dashboard.GET("/getModule1", func(c *gin.Context) { adminDataV1.GetModule1(c) })
			dashboard.POST("/getModule2", func(c *gin.Context) { adminDataV1.GetModule2(c) })
			dashboard.POST("/getModule3", func(c *gin.Context) { adminDataV1.GetModule3(c) })
			dashboard.POST("/getModule4", func(c *gin.Context) { adminDataV1.GetModule4(c) })
		}
	}

	apiRunV1Group := r.Group("/run/v1")
	{
		apiRunV1Group.POST("/get", func(c *gin.Context) { monitorV1.Get(c) })
		apiRunV1Group.GET("/GetGeneralInformation", func(c *gin.Context) { monitorV1.GetGeneralInformation(c) })
		apiRunV1Group.POST("/antModify", func(c *gin.Context) { monitorV1.AntModify(c) })
		apiRunV1Group.POST("/getLeakPart", func(c *gin.Context) { monitorV1.GetLeakPart(c) })
	}

	err = r.Run(fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port")))
	if err != nil {
		return
	}
}
