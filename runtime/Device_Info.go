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
		oi.Device[v.DeviceCode] = models.DeviceInfo{v.ID, v.DeviceCode,
			v.DeviceSequence,
			v.DeviceName,
			v.DeviceInfo}
		oi.Device_Queue = append(oi.Device_Queue, models.DeviceInfo{ID: uint(i + 1), DeviceCode: v.DeviceCode,
			DeviceSequence: v.DeviceSequence, DeviceName: v.DeviceName, DeviceInfo: v.DeviceInfo})
	}
	return true
}
