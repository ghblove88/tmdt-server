package V1

import (
	"TmdtServer/common"
	"TmdtServer/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func Welcome(ginC *gin.Context) {
	ginC.Header("Content-Type", "text/html; charset=utf-8")
	ginC.String(200, `<div style="text-align:center;"><h1>洗消主机后台服务正在运行中……</h1></div>`)
	ginC.String(200, `<div style="text-align:center;"><a href="/admin/v1/system/getSystemInfo">查看系统信息</a></div>`)
}

func Login(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	Username := post["username"].(string)
	Password := post["password"].(string)

	if len(Username) == 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "用户名不能为空"})
		return
	}
	if len(Password) == 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "密码不能为空"})
		return
	}
	engine, err := models.NewOrm()
	if err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "读取数据库出错"})
		return
	}

	var userDb []models.User
	engine.Where("username = ? ", Username).Find(&userDb)

	if len(userDb) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "用户不存在"})
		return
	}

	if strings.EqualFold(userDb[0].Password, common.PwdHash(Password)) {
		mi := common.CacheIntance()
		token := common.Rand().Hex()
		mi[token] = token
		ginC.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"username": Username,
			"roles":        []string{Username, "aaa"},
			"accessToken":  "eyJhbGciOiJIUzUxMiJ9.admin",
			"refreshToken": "eyJhbGciOiJIUzUxMiJ9.adminRefresh",
			"expires":      "2030/10/30 00:00:00"}})
	} else {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "密码错误"})
	}
}

func GetAsyncRoutes(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": []string{}})
}
