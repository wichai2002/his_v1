package domain

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// DateFormat is the standard date format used throughout the application
const DateFormat = "2006-01-02"

// MaxPatientAge is the maximum allowed age for a patient (150 years)
const MaxPatientAge = 150

type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
	Other  Gender = "OTHER"
)

// Scan implements the sql.Scanner interface with proper nil and type handling
func (g *Gender) Scan(value interface{}) error {
	if value == nil {
		*g = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*g = Gender(v)
	case []byte:
		*g = Gender(string(v))
	default:
		return fmt.Errorf("unsupported type for Gender: %T", value)
	}
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

// Scan implements the sql.Scanner interface with proper nil and type handling
func (b *BloodGrp) Scan(value interface{}) error {
	if value == nil {
		*b = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*b = BloodGrp(v)
	case []byte:
		*b = BloodGrp(string(v))
	default:
		return fmt.Errorf("unsupported type for BloodGrp: %T", value)
	}
	return nil
}

func (b BloodGrp) Value() (driver.Value, error) {
	return string(b), nil
}

// Patient model - tenant isolation is handled at schema level
type Patient struct {
	gorm.Model
	FirstNameTH  string    `json:"first_name_th" gorm:"not null,max=255"`
	LastNameTH   string    `json:"last_name_th" gorm:"not null,max=255"`
	MiddleNameTH string    `json:"middle_name_th" gorm:"null,max=255"`
	FirstNameEN  string    `json:"first_name_en" gorm:"not null,max=255"`
	LastNameEN   string    `json:"last_name_en" gorm:"not null,max=255"`
	MiddleNameEN string    `json:"middle_name_en" gorm:"null,max=255"`
	DateOfBirth  time.Time `json:"date_of_birth" gorm:"not null"`
	NickNameTH   string    `json:"nick_name_th" gorm:"not null,max=255"`
	NickNameEN   string    `json:"nick_name_en" gorm:"not null,max=255"`
	PatientHN    string    `json:"patient_hn" gorm:"uniqueIndex;not null"`
	NationalID   string    `json:"national_id" gorm:"uniqueIndex"`
	PassportID   string    `json:"passport_id" gorm:"uniqueIndex"`
	PhoneNumber  string    `json:"phone_number" gorm:"uniqueIndex,max=10"`
	Email        string    `json:"email" gorm:"uniqueIndex, email"`
	// Enum: M, F, OTHER
	Gender      Gender `json:"gender" gorm:"not null,max=5,enum=M|F|OTHER"`
	Nationality string `json:"nationality" gorm:"not null,max=100"`
	// Enum: A, B, O, AB
	BloodGrp BloodGrp `json:"blood_grp" gorm:"not null,max=3,enum=A|B|O|AB"`
}

// DTO for creating a patient
type PatientCreateRequest struct {
	FirstNameTH  string `json:"first_name_th" binding:"required,max=255"`
	LastNameTH   string `json:"last_name_th" binding:"required,max=255"`
	MiddleNameTH string `json:"middle_name_th" binding:"max=255"`
	FirstNameEN  string `json:"first_name_en" binding:"required,max=255"`
	LastNameEN   string `json:"last_name_en" binding:"required,max=255"`
	MiddleNameEN string `json:"middle_name_en" binding:"max=255"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	NickNameTH   string `json:"nick_name_th" binding:"required,max=50"`
	NickNameEN   string `json:"nick_name_en" binding:"required,max=50"`
	NationalID   string `json:"national_id" binding:"required,max=13"`
	PassportID   string `json:"passport_id" binding:"required,max=13"`
	PhoneNumber  string `json:"phone_number" binding:"required,max=10"`
	Email        string `json:"email" binding:"required,email"`
	Gender       Gender `json:"gender" binding:"required,oneof=M F OTHER"`
	Nationality  string `json:"nationality" binding:"required,max=100"`
	BloodGrp     string `json:"blood_grp" binding:"required,oneof=A B O AB"`
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
			dob, err := time.Parse(DateFormat, value.(string))
			if err != nil {
				return nil, fmt.Errorf("invalid date format: %w", err)
			}
			updates[dbTag] = dob
			continue
		}

		updates[dbTag] = value
	}

	return updates, nil
}

// PatientRepository interface - tenant schema provides isolation, no hospitalID needed
type PatientRepository interface {
	GetAll(schemaName string) ([]Patient, error)
	GetByID(id uint, schemaName string) (*Patient, error)
	Search(query string, schemaName string) ([]Patient, error)
	SearchByID(id uint, schemaName string) (*Patient, error)
	Create(patient *Patient, schemaName string) error
	Update(patient *Patient, schemaName string) error
	PartialUpdate(id uint, updates map[string]interface{}, schemaName string) error
	Delete(id uint, schemaName string) error
}

// PatientService interface - tenant isolation handled at schema level
type PatientService interface {
	Search(query string, schemaName string) ([]Patient, error)
	SearchByID(id uint, schemaName string) (*Patient, error)
	Create(req *PatientCreateRequest, schemaName string) (*Patient, error)
	Update(id uint, req *PatientUpdateRequest, schemaName string) (*Patient, error)
	PartialUpdate(id uint, req *PatientPartialUpdateRequest, schemaName string) (*Patient, error)
	Delete(id uint, schemaName string) error
}
