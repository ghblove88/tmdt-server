package runtime

import (
	"EcdsServer/common"
	"EcdsServer/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type RegularEvents struct {
}

func (R *RegularEvents) Start(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		R.checkAndAlertAlarm()
	}
}
func (R *RegularEvents) checkAndAlertAlarm() {
	//检测消毒液是否过期
	R.LiquidExpired()

	// 其他检测项目 ……
}

// LiquidExpired 检测消毒液是否过期
func (R *RegularEvents) LiquidExpired() {
	//清空 记录
	G_Liquid_Expired_Item = G_Liquid_Expired_Item[:0]
	sqlStr := "select a.id,device_name from liquid_device as a  where device_type<> '暂停使用' " +
		"and (IFNULL(validity_period-timestampdiff(day, (select date from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) ,now()),0)<=a.validity_period/3 " +
		"or (select usage_count from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid=a.id ORDER BY id desc limit 0,1) <=a.limit_number/3 ) ORDER BY id"

	var device []map[string]interface{}
	orm, _ := models.NewOrm()
	orm.Raw(sqlStr).Find(&device)
	for _, val := range device {

		sqlStr = fmt.Sprintf("select usage_count, IFNULL(validity_period-timestampdiff(day, date ,now()),0) as remaining  from liquid_detection_record where conclusion='合格' and operator_type='更换' and deviceid='%d' ORDER BY id desc limit 0,1",
			val["id"].(uint32))
		var res []map[string]interface{}
		orm.Raw(sqlStr).Find(&res)

		str := ""
		if len(res) <= 0 {

		} else {
			str = fmt.Sprintf("设备:%s-有效期:%d/有效次数:%d", val["device_name"].(string), res[0]["remaining"].(int64), res[0]["usage_count"].(int64))
			G_Liquid_Expired_Item = append(G_Liquid_Expired_Item, str)
			common.MsgQueueError.Push(str)
		}
	}
	if len(device) > 0 && (G_Liquid_Expired_Count == -1 || time.Now().Minute() <= 5) { //整点语音提醒
		G_Sound_Play.Add_Sound_Queue([]string{"wav/ding.wav", "wav/xdytx.mp3"})
	}
	if len(device) > 0 {
		G_Liquid_Expired_Count = len(device)
	} else {
		G_Liquid_Expired_Count = 0
	}

}
func NewRegularEvents() (*RegularEvents, error) {
	r := &RegularEvents{}
	go r.Start(300) //5分钟检测一次
	zap.S().Info("RegularEvents start")
	return r, nil
}
