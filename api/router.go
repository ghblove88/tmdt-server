package api

import (
	adminDataV1 "TmdtServer/api/admin/V1"
	apiDataV1 "TmdtServer/api/apiData/V1"
	monitorV1 "TmdtServer/api/monitor/V1"
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

	// 创建基于cookie的存储引擎， 参数是用于加密的密钥
	store := cookie.NewStore([]byte("5e59a3e7-2571-4f9b-81ca-1a7682af3aad"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))

	apiDataV1Group := r.Group("/data/v1")
	{

		apiDataV1Group.GET("/operator", func(c *gin.Context) { apiDataV1.GetOperator(c) })
		apiDataV1Group.GET("/device", func(c *gin.Context) { apiDataV1.GetDevice(c) })
		apiDataV1Group.GET("/beds", func(c *gin.Context) { apiDataV1.GetBeds(c) })

	}

	apiAdminV1Group := r.Group("/admin/v1")
	{
		// 服务状态检测提示
		r.GET("/", func(c *gin.Context) { adminDataV1.Welcome(c) })

		apiAdminV1Group.POST("/login", func(c *gin.Context) { adminDataV1.Login(c) })
		apiAdminV1Group.GET("/get-async-routes", func(c *gin.Context) { adminDataV1.GetAsyncRoutes(c) })
		apiAdminV1Group.GET("/operator/list", func(c *gin.Context) { adminDataV1.GetOperatorList(c) })

		daily := apiAdminV1Group.Group("/daily")
		{
			daily.POST("/queryDevice", func(c *gin.Context) { adminDataV1.QueryDeviceList(c) })
			daily.POST("/deleteDevice", func(c *gin.Context) { adminDataV1.DeleteDevice(c) })
			daily.POST("/updateDevice", func(c *gin.Context) { adminDataV1.ModifyDevice(c) })

			daily.POST("/queryOperator", func(c *gin.Context) { adminDataV1.QueryOperatorList(c) })
			daily.POST("/deleteOperator", func(c *gin.Context) { adminDataV1.DeleteOperator(c) })
			daily.POST("/updateOperator", func(c *gin.Context) { adminDataV1.ModifyOperator(c) })

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

	}

	apiRunV1Group := r.Group("/run/v1")
	{
		apiRunV1Group.POST("/get", func(c *gin.Context) { monitorV1.Get(c) })
		apiRunV1Group.GET("/GetGeneralInformation", func(c *gin.Context) { monitorV1.GetGeneralInformation(c) })

	}

	err = r.Run(fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port")))
	if err != nil {
		return
	}
}
