package runtime

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

// Bed represents the bed table in the database.
type Bed struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	BedCode string `gorm:"unique;not null"`
	BedInfo string `gorm:"not null"`
}

// BedService provides methods to operate on the bed table.
type BedService struct {
	db    *gorm.DB
	cache map[uint]Bed
	mutex sync.RWMutex
}

// NewBedService creates a new BedService and loads all bed records into the cache.
func NewBedService(db *gorm.DB) (*BedService, error) {
	bedService := &BedService{
		db:    db,
		cache: make(map[uint]Bed),
	}
	if err := bedService.loadCache(); err != nil {
		return nil, err
	}
	return bedService, nil
}

// loadCache loads all bed records into the cache.
func (s *BedService) loadCache() error {
	var beds []Bed
	if err := s.db.Find(&beds).Error; err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, bed := range beds {
		s.cache[bed.ID] = bed
	}
	return nil
}

// CreateBed adds a new bed record to the database and updates the cache.
func (s *BedService) CreateBed(bedCode, bedInfo string) (Bed, error) {
	bed := Bed{BedCode: bedCode, BedInfo: bedInfo}
	if err := s.db.Create(bed).Error; err != nil {
		return Bed{BedCode: "", BedInfo: ""}, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache[bed.ID] = bed
	return bed, nil
}

// UpdateBed updates an existing bed record in the database and the cache.
func (s *BedService) UpdateBed(id uint, bedCode, bedInfo string) (Bed, error) {
	bed := Bed{}
	if err := s.db.First(bed, id).Error; err != nil {
		return Bed{BedCode: "", BedInfo: ""}, err
	}
	bed.BedCode = bedCode
	bed.BedInfo = bedInfo
	if err := s.db.Save(bed).Error; err != nil {
		return Bed{BedCode: "", BedInfo: ""}, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache[bed.ID] = bed
	return bed, nil
}

// DeleteBed deletes a bed record from the database and updates the cache.
func (s *BedService) DeleteBed(id uint) error {
	if err := s.db.Delete(&Bed{}, id).Error; err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.cache, id)
	return nil
}

// GetBed retrieves a bed record from the cache.
func (s *BedService) GetBed(id uint) (Bed, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if bed, exists := s.cache[id]; exists {
		return bed, nil
	}
	return Bed{BedCode: "", BedInfo: ""}, fmt.Errorf("bed with ID %d not found", id)
}

// GetAllBeds retrieves all bed records from the cache.
func (s *BedService) GetAllBeds() ([]Bed, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	beds := make([]Bed, 0, len(s.cache))
	for _, bed := range s.cache {
		beds = append(beds, bed)
	}
	return beds, nil
}
