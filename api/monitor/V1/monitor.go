package V1

import (
	"TmdtServer/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Get(ginC *gin.Context) {
	//测试
	if *common.TestMode {

	}

	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "有数据更新", "hash": "hash",
		"antList": "antList", "datetime": time.Now().Format("2006-01-02 15:04:05")})
}

func GetGeneralInformation(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功",
		"error":   common.MsgQueueError.PopOrWait(1),
		"warning": common.MsgQueueWarning.PopOrWait(1),
		"info":    common.MsgQueueInfo.PopOrWait(1)})
}
