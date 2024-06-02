package runtime

import (
	models "TmdtServer/models"
)

type Device_Info struct {
	Device       map[string]models.DeviceInfo
	Device_Queue []models.DeviceInfo
}

func (oi *Device_Info) GetAll() (res []models.DeviceInfo) {
	return oi.Device_Queue
}

func (oi *Device_Info) GetAt(number string) (res models.DeviceInfo) {
	return oi.Device[number]
}

func (oi *Device_Info) PullDB() (res bool) {
	oi.Device = make(map[string]models.DeviceInfo)
	oi.Device_Queue = []models.DeviceInfo{}

	engine, err := models.NewOrm()
	if err != nil {
		println(err.Error())
		return
	}

	var device []models.DeviceInfo
	db := engine.Where("id >?", 0).Order("id").Find(&device)
	if db.Error != nil {
		return
	}

	for i, v := range device {
		oi.Device[v.EndoscopeNumber] = models.DeviceInfo{v.Id, v.EndoscopeNumber,
			v.EndoscopeType,
			v.EndoscopeInfo,
			v.Status}
		oi.Device_Queue = append(oi.Device_Queue, models.DeviceInfo{Id: i + 1, EndoscopeNumber: v.EndoscopeNumber,
			EndoscopeType: v.EndoscopeType, EndoscopeInfo: v.EndoscopeInfo, Status: v.Status})
	}
	return true
}
