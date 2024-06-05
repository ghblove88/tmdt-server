package V1

import (
	"TmdtServer/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOperatorList(ginC *gin.Context) {
	res, err := models.Operatorlist()
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "data": err.Error()})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
