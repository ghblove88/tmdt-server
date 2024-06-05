package models

import "time"

// TemperatureRecord represents the temperature_records table in the database.
type TemperatureRecord struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	DeviceCode   string  `gorm:"not null"`
	BedCode      string  `gorm:"not null"`
	Operator     string  `gorm:"not null"`
	Patient      string  `gorm:"not null"`
	Temperature1 float64 `gorm:"not null"`
	Temperature2 float64
	Temperature3 float64
	RecordTime   time.Time `gorm:"not null"`
}

// DeviceInfo represents the device table in the database.
type DeviceInfo struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	DeviceCode     string `gorm:"unique;not null"`
	DeviceSequence string `gorm:"unique;not null"`
	DeviceName     string `gorm:"not null"`
	DeviceInfo     string `gorm:"not null"`
}

// Bed represents the bed table in the database.
type Bed struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	BedCode string `gorm:"unique;not null"`
	BedInfo string `gorm:"not null"`
}

type OperatorInfo struct {
	Id     int    `xorm:"not null pk autoincr INT"`
	Number string `xorm:"VARCHAR(100)"`
	Name   string `xorm:"VARCHAR(100)"`
}

type User struct {
	Id       int    `xorm:"not null pk autoincr INT"`
	Username string `xorm:"VARCHAR(100)"`
	Password string `xorm:"CHAR(36)"`
}
