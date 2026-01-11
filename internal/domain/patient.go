package domain

import (
	"database/sql/driver"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
	Other  Gender = "OTHER"
)

func (g *Gender) Scan(value interface{}) error {
	*g = Gender(value.(string))
	return nil
}
func (g Gender) Value() (driver.Value, error) {
	return string(g), nil
}

type BloodGrp string

const (
	A  BloodGrp = "A"
	B  BloodGrp = "B"
	O  BloodGrp = "O"
	AB BloodGrp = "AB"
)

func (b *BloodGrp) Scan(value interface{}) error {
	*b = BloodGrp(value.(string))
	return nil
}

func (b BloodGrp) Value() (driver.Value, error) {
	return string(b), nil
}

// Patient model with foreign key to Hospital
type Patient struct {
	gorm.Model
	FirstNameTH  string    `json:"first_name_th" gorm:"not null, max=255"`
	LastNameTH   string    `json:"last_name_th" gorm:"not null, max=255"`
	MiddleNameTH string    `json:"middle_name_th" gorm:"not null, max=255"`
	FirstNameEN  string    `json:"first_name_en" gorm:"not null, max=255"`
	LastNameEN   string    `json:"last_name_en" gorm:"not null, max=255"`
	MiddleNameEN string    `json:"middle_name_en" gorm:"not null, max=255"`
	DateOfBirth  time.Time `json:"date_of_birth" gorm:"not null"`
	NickNameTH   string    `json:"nick_name_th" gorm:"not null, max=255"`
	NickNameEN   string    `json:"nick_name_en" gorm:"not null, max=255"`
	PatientHN    string    `json:"patient_hn" gorm:"uniqueIndex;not null"`
	NationalID   string    `json:"national_id" gorm:"uniqueIndex"`
	PassportID   string    `json:"passport_id" gorm:"uniqueIndex"`
	PhoneNumber  string    `json:"phone_number" gorm:"uniqueIndex, max=10"`
	Email        string    `json:"email" gorm:"uniqueIndex, email"`
	// Enum: M, F, OTHER
	Gender      Gender `json:"gender" gorm:"not null,max=5,enum=M|F|OTHER"`
	Nationality string `json:"nationality" gorm:"not null, max=100"`
	// Enum: A, B, O, AB
	BloodGrp   BloodGrp `json:"blood_grp" gorm:"not null,max=3,enum=A|B|O|AB"`
	HospitalID uint     `json:"hospital_id" gorm:"not null"`
	Hospital   Hospital `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"`
}

// DTO for creating a patient
type PatientCreateRequest struct {
	FirstNameTH  string `json:"first_name_th" binding:"required, max=255"`
	LastNameTH   string `json:"last_name_th" binding:"required, max=255"`
	MiddleNameTH string `json:"middle_name_th" binding:"required, max=255"`
	FirstNameEN  string `json:"first_name_en" binding:"required, max=255"`
	LastNameEN   string `json:"last_name_en" binding:"required, max=255"`
	MiddleNameEN string `json:"middle_name_en" binding:"required, max=255"`
	DateOfBirth  string `json:"date_of_birth" binding:"required, date_format=2006-01-02"`
	NickNameTH   string `json:"nick_name_th" binding:"required, max=50"`
	NickNameEN   string `json:"nick_name_en" binding:"required, max=50"`
	NationalID   string `json:"national_id" binding:"required, max=13"`
	PassportID   string `json:"passport_id" binding:"required, max=13"`
	PhoneNumber  string `json:"phone_number" binding:"required, max=10"`
	Email        string `json:"email" binding:"required, email"`
	Gender       Gender `json:"gender" binding:"required, enum=M|F|OTHER"`
	Nationality  string `json:"nationality" binding:"required, max=100"`
	BloodGrp     string `json:"blood_grp" binding:"required, enum=A|B|O|AB"`
	HospitalID   uint   `json:"hospital_id" binding:"required"`
}

// DTO for full update (PUT) - all fields are required
type PatientUpdateRequest struct {
	FirstNameTH  string `json:"first_name_th" binding:"required,max=255"`
	LastNameTH   string `json:"last_name_th" binding:"required,max=255"`
	MiddleNameTH string `json:"middle_name_th" binding:"max=255"`
	FirstNameEN  string `json:"first_name_en" binding:"required,max=255"`
	LastNameEN   string `json:"last_name_en" binding:"required,max=255"`
	MiddleNameEN string `json:"middle_name_en" binding:"max=255"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	NickNameTH   string `json:"nick_name_th" binding:"max=50"`
	NickNameEN   string `json:"nick_name_en" binding:"max=50"`
	NationalID   string `json:"national_id" binding:"max=13"`
	PassportID   string `json:"passport_id" binding:"max=13"`
	PhoneNumber  string `json:"phone_number" binding:"max=10"`
	Email        string `json:"email" binding:"omitempty,email"`
	Gender       string `json:"gender" binding:"required,oneof=M F OTHER"`
	Nationality  string `json:"nationality" binding:"required,max=100"`
	BloodGrp     string `json:"blood_grp" binding:"required,oneof=A B O AB"`
}

// DTO for partial update (PATCH) - uses pointers to distinguish between nil and empty values
type PatientPartialUpdateRequest struct {
	FirstNameTH  *string `json:"first_name_th,omitempty" binding:"omitempty,max=255" db:"first_name_th"`
	LastNameTH   *string `json:"last_name_th,omitempty" binding:"omitempty,max=255" db:"last_name_th"`
	MiddleNameTH *string `json:"middle_name_th,omitempty" binding:"omitempty,max=255" db:"middle_name_th"`
	FirstNameEN  *string `json:"first_name_en,omitempty" binding:"omitempty,max=255" db:"first_name_en"`
	LastNameEN   *string `json:"last_name_en,omitempty" binding:"omitempty,max=255" db:"last_name_en"`
	MiddleNameEN *string `json:"middle_name_en,omitempty" binding:"omitempty,max=255" db:"middle_name_en"`
	DateOfBirth  *string `json:"date_of_birth,omitempty" db:"date_of_birth"`
	NickNameTH   *string `json:"nick_name_th,omitempty" binding:"omitempty,max=50" db:"nick_name_th"`
	NickNameEN   *string `json:"nick_name_en,omitempty" binding:"omitempty,max=50" db:"nick_name_en"`
	NationalID   *string `json:"national_id,omitempty" binding:"omitempty,max=13" db:"national_id"`
	PassportID   *string `json:"passport_id,omitempty" binding:"omitempty,max=13" db:"passport_id"`
	PhoneNumber  *string `json:"phone_number,omitempty" binding:"omitempty,max=10" db:"phone_number"`
	Email        *string `json:"email,omitempty" binding:"omitempty,email" db:"email"`
	Gender       *string `json:"gender,omitempty" binding:"omitempty,oneof=M F OTHER" db:"gender"`
	Nationality  *string `json:"nationality,omitempty" binding:"omitempty,max=100" db:"nationality"`
	BloodGrp     *string `json:"blood_grp,omitempty" binding:"omitempty,oneof=A B O AB" db:"blood_grp"`
}

// ToMap converts non-nil fields to a map for partial updates
func (r *PatientPartialUpdateRequest) ToMap() (map[string]interface{}, error) {
	updates := make(map[string]interface{})
	v := reflect.ValueOf(r).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsNil() {
			continue
		}

		dbTag := t.Field(i).Tag.Get("db")
		value := field.Elem().Interface()

		// Handle date_of_birth parsing
		if dbTag == "date_of_birth" {
			dob, err := time.Parse("2006-01-02", value.(string))
			if err != nil {
				return nil, err
			}
			updates[dbTag] = dob
			continue
		}

		updates[dbTag] = value
	}

	return updates, nil
}

type PatientRepository interface {
	GetAll() ([]Patient, error)
	GetByID(id uint) (*Patient, error)
	Search(query string, hospitalID uint) ([]Patient, error)
	SearchByID(id uint, hospitalID uint) (*Patient, error)
	Create(patient *Patient) error
	Update(patient *Patient) error
	PartialUpdate(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

type PatientService interface {
	Search(query string, hospitalID uint) ([]Patient, error)
	SearchByID(id uint, hospitalID uint) (*Patient, error)
	Create(req *PatientCreateRequest, hospitalID uint) (*Patient, error)
	Update(id uint, req *PatientUpdateRequest) (*Patient, error)
	PartialUpdate(id uint, req *PatientPartialUpdateRequest) (*Patient, error)
	Delete(id uint) error
}
