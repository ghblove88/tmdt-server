package models

// DeviceTypelist 获取设备类型列表
// 参数: 无
// 返回值:
// rows []DeviceType - 设备类型列表
// res error - 操作过程中可能出现的错误
func DeviceTypelist() (rows []DeviceType, res error) {

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

// DeviceTypeDelete 根据提供的ID列表删除设备类型记录。
// ids: 表示需要删除的设备类型ID列表，为整数数组。
// res: 返回操作结果，如果删除成功，则返回nil；如果删除过程中出现错误，则返回相应的错误信息。
func DeviceTypeDelete(ids []int) (res error) {
	// 新建ORM对象以进行数据库操作
	if o, err := NewOrm(); err != nil {
		// 如果ORM对象创建失败，则直接返回错误
		return
	} else {
		// 使用ORM对象删除指定ID的设备类型记录
		o.Delete(&DeviceType{}, ids)
	}
	return
}
