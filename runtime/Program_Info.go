package runtime

import (
	"EcdsServer/models"
	"fmt"
	"strconv"
)

type STRUCT_PROGRAM_INFO struct {
	Name          string                              `json:"name"`
	TotalCostTime int                                 `json:"TotalCostTime"`
	StepList      map[string]STRUCT_PROGRAM_LIST_INFO `json:"StepList"`
}

type STRUCT_PROGRAM_LIST_INFO struct {
	Step      string `json:"Step"`
	Cost_time int    `json:"Step_time"`
}

type Program_Info struct {
	Program        map[string]STRUCT_PROGRAM_INFO
	ProgramOrderly map[int]STRUCT_PROGRAM_INFO
}

func (oi *Program_Info) GetAll() (res map[int]STRUCT_PROGRAM_INFO) {
	return oi.ProgramOrderly
}

func (oi *Program_Info) PullDB() bool {
	oi.Program = make(map[string]STRUCT_PROGRAM_INFO)
	oi.ProgramOrderly = make(map[int]STRUCT_PROGRAM_INFO)
	engine, err := models.NewOrm()
	if err != nil {
		println(err.Error())
		return false
	}
	var program []models.Program
	db := engine.Where("id >?", 0).Order("id").Find(&program)
	if db.Error != nil {

		return false
	}

	for i, v := range program {
		step_list := make(map[string]STRUCT_PROGRAM_LIST_INFO)
		var programList []models.ProgramList
		db := engine.Where("name =?", v.Name).Order("id").Find(&programList)
		if db.Error != nil {
			fmt.Println(db.Error)
			return false
		}

		for si, sv := range programList {

			step_list[strconv.Itoa(si+1)] = STRUCT_PROGRAM_LIST_INFO{sv.Step,
				sv.CostTime}
			si++
		}

		oi.Program[v.Name] = STRUCT_PROGRAM_INFO{v.Name,
			v.TotalCostTime,
			step_list}

		oi.ProgramOrderly[i+1] = STRUCT_PROGRAM_INFO{v.Name,
			v.TotalCostTime,
			step_list}
	}
	return true
}
