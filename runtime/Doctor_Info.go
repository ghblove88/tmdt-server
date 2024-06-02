package runtime

import (
	"TmdtServer/models"
	"strconv"
)

type Doctor_Info struct {
	Doctor      map[string]models.DoctorInfo
	DoctorQueue []models.DoctorInfo
}

func (oi *Doctor_Info) GetAll() (res []models.DoctorInfo) {
	return oi.DoctorQueue
}

func (oi *Doctor_Info) PullDB() (res bool) {

	oi.Doctor = map[string]models.DoctorInfo{}
	oi.DoctorQueue = []models.DoctorInfo{}

	engine, err := models.NewOrm()
	if err != nil {
		println(err.Error())
		return
	}

	var doctor []models.DoctorInfo
	db := engine.Where("id >?", 0).Order("id").Find(&doctor)
	if db.Error != nil {
		return
	}

	for i, v := range doctor {
		oi.Doctor[strconv.Itoa(v.Id)] = models.DoctorInfo{i, v.Name}
		oi.DoctorQueue = append(oi.DoctorQueue, models.DoctorInfo{Id: i + 1, Name: v.Name})
	}
	return true
}
