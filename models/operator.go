package models

func Operatorlist() (rows []OperatorInfo, res error) {

	// 新建Orm对象用于数据库操作
	if o, err := NewOrm(); err != nil {
		// 如果创建Orm对象失败，直接返回错误
		return
	} else {
		// 使用Orm对象查询设备类型列表，并按id排序
		o.Find(&rows).Order("id")
	}
	return rows, nil
}
