package runtime

import (
	"gorm.io/gorm"
	"time"
)

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

// TemperatureRecordService provides methods to operate on the temperature_records table.
type TemperatureRecordService struct {
	db *gorm.DB
}

// NewTemperatureRecordService creates a new TemperatureRecordService.
func NewTemperatureRecordService(db *gorm.DB) *TemperatureRecordService {
	return &TemperatureRecordService{db: db}
}

// CreateTemperatureRecord adds a new temperature record to the database.
func (s *TemperatureRecordService) CreateTemperatureRecord(deviceCode, bedCode, operator, patient string, temp1, temp2, temp3 float64, recordTime time.Time) (*TemperatureRecord, error) {
	record := &TemperatureRecord{
		DeviceCode:   deviceCode,
		BedCode:      bedCode,
		Operator:     operator,
		Patient:      patient,
		Temperature1: temp1,
		Temperature2: temp2,
		Temperature3: temp3,
		RecordTime:   recordTime,
	}
	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// UpdateTemperatureRecord updates an existing temperature record.
func (s *TemperatureRecordService) UpdateTemperatureRecord(id uint, deviceCode, bedCode, operator, patient string, temp1, temp2, temp3 float64, recordTime time.Time) (*TemperatureRecord, error) {
	record := &TemperatureRecord{}
	if err := s.db.First(record, id).Error; err != nil {
		return nil, err
	}
	record.DeviceCode = deviceCode
	record.BedCode = bedCode
	record.Operator = operator
	record.Patient = patient
	record.Temperature1 = temp1
	record.Temperature2 = temp2
	record.Temperature3 = temp3
	record.RecordTime = recordTime
	if err := s.db.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// DeleteTemperatureRecord deletes a temperature record from the database.
func (s *TemperatureRecordService) DeleteTemperatureRecord(id uint) error {
	if err := s.db.Delete(&TemperatureRecord{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetTemperatureRecord retrieves a temperature record by its ID.
func (s *TemperatureRecordService) GetTemperatureRecord(id uint) (*TemperatureRecord, error) {
	record := &TemperatureRecord{}
	if err := s.db.First(record, id).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// GetAllTemperatureRecords retrieves all temperature records from the database.
func (s *TemperatureRecordService) GetAllTemperatureRecords() ([]TemperatureRecord, error) {
	var records []TemperatureRecord
	if err := s.db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
