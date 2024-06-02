package runtime

import (
	"EcdsServer/models"
	"strconv"
)

type STRUCT_TIMEPLAN_INFO struct {
	Name               string `json:"name"`
	Rinse_time         int    `json:"rinse_time"`
	Rinse_solution     string `json:"rinse_solution"`
	Disinfect_time     int    `json:"disinfect_time"`
	Sterilise_time     int    `json:"sterilise_time"`
	Disinfect_solution string `json:"disinfect_solution"`
}

type TimePlan_Info struct {
	Timeplan map[string]STRUCT_TIMEPLAN_INFO
}

func (oi *TimePlan_Info) GetAll() (res map[string]STRUCT_TIMEPLAN_INFO) {
	res = make(map[string]STRUCT_TIMEPLAN_INFO)
	i := 0
	for _, value := range oi.Timeplan {
		res[strconv.Itoa(i+1)] = value
		i++
	}
	return res
}

func (oi *TimePlan_Info) GetAt(device_Type string, Choose int) (timelong int) {
	res := oi.Timeplan[device_Type]
	if Choose == 0 { // 酶洗液
		return res.Rinse_time
	} else if Choose == 1 { // 消毒液
		return res.Disinfect_time
	} else if Choose == 2 { // 灭菌
		return res.Sterilise_time
	}
	return 0
}

func (oi *TimePlan_Info) PullDB() (res bool) {

	oi.Timeplan = make(map[string]STRUCT_TIMEPLAN_INFO)

	engine, err := models.NewOrm()
	if err != nil {
		println(err.Error())
		return
	}

	var timeplan []models.TimePlan
	engine.Where("id >?", 0).Order("id").Find(&timeplan)
	if err != nil {
		return
	}
	for _, v := range timeplan {
		oi.Timeplan[v.Name] = STRUCT_TIMEPLAN_INFO{v.Name,
			v.RinseTime,
			v.RinseSolution,
			v.DisinfectTime,
			v.SteriliseTime,
			v.DisinfectSolution}
	}
	return true
}
