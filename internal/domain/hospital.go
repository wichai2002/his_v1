package domain

import (
	"gorm.io/gorm"
)

type Hospital struct {
	gorm.Model
	Name               string `json:"name" gorm:"uniqueIndex;not null, max=255"`
	HospitalCode       string `json:"hospital_code" gorm:"uniqueIndex;not null, max=10"`
	PhoneNumber        string `json:"phone_number" gorm:"uniqueIndex;not null, max=10"`
	Email              string `json:"email" gorm:"uniqueIndex;not null, email"`
	Address            string `json:"address" gorm:"not null, type:text"`
	HNRunningNumber    int    `json:"hn_running_number" gorm:"default:0"`
	StaffRunningNumber int    `json:"staff_running_number" gorm:"default:0"`
}

type HospitalPublicResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type HospitalRepository interface {
	GetAllPublic() ([]HospitalPublicResponse, error)
	GetAll() ([]Hospital, error)
	GetByID(id uint) (*Hospital, error)
	Create(hospital *Hospital) error
	Update(hospital *Hospital) error
	Delete(id uint) error
	IncrementHNRunningNumber(id uint) (int, error)
	IncrementStaffRunningNumber(id uint) (int, error)
}

type HospitalService interface {
	GetAllPublic() ([]HospitalPublicResponse, error)
	GetAll() ([]Hospital, error)
	GetByID(id uint) (*Hospital, error)
	Create(hospital *Hospital) error
	Update(hospital *Hospital) error
	Delete(id uint) error
	GenerateHN(hospitalID uint) (string, error)
	GenerateStaffCode(hospitalID uint) (string, error)
}
