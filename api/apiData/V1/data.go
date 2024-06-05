package V1

import (
	"TmdtServer/runtime"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOperator(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Operator_Info.GetAll())
}

func GetDevice(ginC *gin.Context) {
	ginC.JSON(http.StatusOK, runtime.G_Device_Info.GetAll())
}

func GetBeds(ginC *gin.Context) {
	beds, err := runtime.G_BedService.GetAllBeds()
	if err != nil {
		return
	}
	ginC.JSON(http.StatusOK, beds)
}
