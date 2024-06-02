package V1

import (
	"TmdtServer/common"
	"TmdtServer/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func BackupData(ginC *gin.Context) {
	common.TermExec("rm", []string{"-rf", "./bgStatic/dump.sql.gz"})
	err := CreateDump("./bgStatic/")
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "备份失败:CreateDump" + err.Error()})
		return
	}
	cmd := exec.Command("gzip", []string{"./bgStatic/dump.sql"}...)
	if err := cmd.Start(); err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "备份失败:gzip" + err.Error()})
		log.Println("Start: ", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "备份失败:Wait" + err.Error()})
		log.Println("Wait: ", err.Error())
		return
	}
	//err = cmd.Run()
	//if err != nil {
	//	ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "备份失败" + err.Error()})
	//	return
	//}
	//add := fmt.Sprintf("%s:%d", viper.GetString("web_server.address"), viper.GetInt("web_server.port"))
	ginC.JSON(http.StatusOK, gin.H{"success": true, "data": "/bgStatic/dump.sql.gz"})

}
func CreateDump(dumpDir string) error {

	var args []string
	dbUser := common.Config.GetString("ecds_server.username")
	dbPass := common.Config.GetString("ecds_server.password")
	dbAddress := common.Config.GetString("ecds_server.address")
	dbPort := common.Config.GetString("ecds_server.port")
	dbName := common.Config.GetString("ecds_server.name")

	args = append(args, "-u"+dbUser)
	args = append(args, "-p"+dbPass)
	args = append(args, "-h"+dbAddress)
	args = append(args, "-P"+dbPort)
	args = append(args, "-e")
	args = append(args, "--max_allowed_packet=1048576")
	args = append(args, "--net_buffer_length=16384")
	args = append(args, dbName)

	out, err := os.OpenFile(path.Join(dumpDir, "dump.sql"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("openFile: ", err.Error())
		return err
	}
	cmd := exec.Command("mysqldump", args...)
	cmd.Stdout = out

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		return err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Start: ", err.Error())
		return err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		fmt.Println("ReadAll stderr: ", err.Error())
		return err
	}

	if len(bytesErr) != 0 {
		fmt.Printf("stderr is not nil: %s", bytesErr)
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: ", err.Error())
		return err
	}
	return cmd.Run()
}

func BackupFileUpload(ginC *gin.Context) {
	f, err := ginC.FormFile("file")
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "上传失败!"})
		return
	}
	pathFile := "./bgStatic/" + f.Filename
	err = ginC.SaveUploadedFile(f, pathFile)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "上传失败!"})
		return
	}
	err = common.TermExec("rm", []string{"-rf", "./bgStatic/dump.sql"})
	err = common.TermExec("gunzip", []string{pathFile})
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "解压缩失败!"})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "上传成功!"})
}

func DataRecovery(ginC *gin.Context) {
	if common.Exist("./bgStatic/dump.sql") == false {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "文件不存在!"})
		return
	}

	err := Dump("./bgStatic/dump.sql")
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "恢复失败!"})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "恢复成功!"})
}

func Dump(dumpDir string) error {
	var args []string
	dbUser := common.Config.GetString("ecds_server.username")
	dbPass := common.Config.GetString("ecds_server.password")
	dbAddress := common.Config.GetString("ecds_server.address")
	dbPort := common.Config.GetString("ecds_server.port")
	dbName := common.Config.GetString("ecds_server.name")

	args = append(args, "-u"+dbUser)
	args = append(args, "-p"+dbPass)
	args = append(args, "-h"+dbAddress)
	args = append(args, "-P"+dbPort)
	args = append(args, dbName)

	in, err := os.OpenFile(dumpDir, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("openfile: ", err.Error())
		return err
	}
	cmd := exec.Command("mysql", args...)
	cmd.Stdin = in

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		return err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Start: ", err.Error())
		return err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		fmt.Println("ReadAll stderr: ", err.Error())
		return err
	}

	if len(bytesErr) != 0 {
		fmt.Printf("stderr is not nil: %s", bytesErr)
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: ", err.Error())
		return err
	}
	return cmd.Run()
}

type UserCondition struct {
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

func QueryUserList(ginC *gin.Context) {
	var uc UserCondition
	err := ginC.ShouldBind(&uc)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	conditionStr := " where id>0 "

	if len(uc.Name) > 0 {
		conditionStr += " and username like '%" + uc.Name + "%' "
	}

	var offset int
	if uc.Pagination.CurrentPage <= 1 {
		offset = 0
	} else {
		offset = (uc.Pagination.CurrentPage - 1) * uc.Pagination.PageSize
	}

	orm, err := models.NewOrm()
	if err != nil {
		return
	}

	var total int64
	orm.Raw("select id from user " + conditionStr).Count(&total)

	listSqlStr := "select * from user "

	// 分页排序
	ordLimitStr := " order by " + uc.Sort + " " + uc.Order + " limit " + strconv.FormatInt(int64(offset), 10) + "," + strconv.FormatInt(int64(uc.Pagination.PageSize), 10)

	var rows []map[string]interface{}
	db := orm.Raw(listSqlStr + conditionStr + ordLimitStr).Scan(&rows)
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "用户列表", "data": rows, "Pagination": gin.H{"total": total,
		"currentPage": uc.Pagination.CurrentPage, "pageSize": uc.Pagination.PageSize}})
}

func DeleteUser(ginC *gin.Context) {
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
	orm.Delete(&models.User{}, post[0]["id"])
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}

func UserModify(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	db := orm.Save(&models.User{Username: post["Username"].(string), Password: common.PwdHash("admin123")})
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建成功,默认密码：admin123"})
}
func ResetPassword(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	orm, _ := models.NewOrm()
	db := orm.Model(&models.User{}).Where("id = ?", post["id"]).Update("password", common.PwdHash("admin123"))
	if db.Error != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": db.Error})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "重置密码成功,默认密码：admin123"})
}
func ChangePassword(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	engine, err := models.NewOrm()
	if err != nil {
		log.Print(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "读取数据库出错"})
		return
	}
	userName := post["userName"].(string)
	oldPassword := post["oldPassword"].(string)
	newPassword := post["newPassword"].(string)
	confirmPassword := post["confirmPassword"].(string)

	var userDb []models.User
	engine.Where("username = ? ", userName).Find(&userDb)

	if len(userDb) <= 0 {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "用户不存在"})
		return
	}

	if strings.EqualFold(userDb[0].Password, common.PwdHash(oldPassword)) {
		if strings.EqualFold(newPassword, confirmPassword) {
			engine.Model(&models.User{}).Where("username = ? ", userName).Update("password", common.PwdHash(newPassword))
			ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "修改密码成功"})
		} else {
			ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "两次密码输入不一致"})
		}
	} else {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "旧密码错误"})
	}
}
func GetSystemInfo(ginC *gin.Context) {
	//CPU使用率
	cpuPercent, _ := cpu.Percent(time.Second, true)
	//CPU核心数
	cpuNumber, _ := cpu.Counts(true)
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get memory info fail. err： ", err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取内存信息失败"})
		return
	}
	// 获取总内存大小，单位GB
	memTotal := memInfo.Total / 1024 / 1024 / 1024
	// 获取已用内存大小，单位MB
	memUsed := memInfo.Used / 1024 / 1024
	// 可用内存大小
	memAva := memInfo.Available / 1024 / 1024
	// 内存可用率
	memUsedPercent := memInfo.UsedPercent
	//系统平均负载
	loadInfo, err := load.Avg()
	//主机信息
	hostInfo, err := host.Info()

	//获取硬盘存储
	diskPart, err := disk.Partitions(false)
	if err != nil {
		fmt.Println(err)
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取硬盘存储失败"})
		return
	}
	var diskUsage []disk.UsageStat
	for _, dp := range diskPart {
		diskUsed, _ := disk.Usage(dp.Mountpoint)
		diskUsage = append(diskUsage, *diskUsed)
		//fmt.Printf("分区总大小: %d MB \n", diskUsed.Total/1024/1024)
		//fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
		//fmt.Printf("分区inode使用率: %.3f %% \n", diskUsed.InodesUsedPercent)
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取系统信息成功", "data": gin.H{
		"cpuPercent":     cpuPercent,
		"cpuNumber":      cpuNumber,
		"memTotal":       memTotal,
		"memUsed":        memUsed,
		"memAva":         memAva,
		"memUsedPercent": memUsedPercent,
		"loadInfo":       loadInfo,
		"hostInfo":       hostInfo,
		"diskPart":       diskPart,
		"diskUsage":      diskUsage,
	}})

}
func Restart(ginC *gin.Context) {
	fmt.Println("重启程序...")

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径:", err)
		return
	}

	// 使用 syscall.Exec 来替换当前进程
	err = syscall.Exec(exePath, os.Args, os.Environ())
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "重启动失败"})
		fmt.Println("重启失败:", err)
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "重启动成功"})
}

func GetConfigurationFile(ginC *gin.Context) {
	// 指定要读取的文件路径
	filePath := "./app.yaml"

	// 使用 os.ReadFile 函数读取文件内容
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "读取配置文件失败"})
		fmt.Printf("Error reading file '%s': %v\n", filePath, err)
		return
	}

	// 将读取到的字节切片转换为字符串
	contentStr := string(contentBytes)

	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取配置文件成功", "data": contentStr})
}
func SaveConfigurationFile(ginC *gin.Context) {
	var post map[string]interface{}
	err := ginC.BindJSON(&post)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	// 定义要写入的文件路径
	filePath := "./app.yaml"

	// 使用 os.WriteFile 函数写入文件
	err = os.WriteFile(filePath, []byte(post["content"].(string)), 0644)
	if err != nil {
		ginC.JSON(http.StatusOK, gin.H{"success": false, "msg": "写入配置文件失败"})
		return
	}
	ginC.JSON(http.StatusOK, gin.H{"success": true, "msg": "写入配置文件成功"})
}
