package runtime

import (
	"EcdsServer/models"
)

type STRUCT_OPERATOR_INFO struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type Operator_Info struct {
	Operator      map[string]STRUCT_OPERATOR_INFO
	OperatorQueue []models.OperatorInfo
}

func (oi *Operator_Info) GetAll() (res []models.OperatorInfo) {
	return oi.OperatorQueue
}

func (oi *Operator_Info) GetAt(number string) (res STRUCT_OPERATOR_INFO) {
	return oi.Operator[number]
}

func (oi *Operator_Info) ReverseAt(name string) (res string) {
	number := ""
	for _, value := range oi.Operator {
		if value.Name == name {
			number = value.Number
		}
	}
	return number
}

func (oi *Operator_Info) PullDB() (res bool) {

	oi.Operator = make(map[string]STRUCT_OPERATOR_INFO)
	oi.OperatorQueue = []models.OperatorInfo{}
	engine, err := models.NewOrm()
	if err != nil {
		println(err.Error())
		return
	}

	var operator []models.OperatorInfo
	db := engine.Where("id >?", 0).Order("id").Find(&operator)
	if db.Error != nil {
		return
	}

	for i, v := range operator {
		oi.Operator[v.Number] = STRUCT_OPERATOR_INFO{v.Number, v.Name}
		oi.OperatorQueue = append(oi.OperatorQueue, models.OperatorInfo{Id: i + 1, Number: v.Number, Name: v.Name})
	}
	return true
}
