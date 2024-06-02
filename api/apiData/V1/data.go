package V1

import (
	"TmdtServer/runtime"
	"github.com/gin-gonic/gin"
	"net/http"
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
