package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wichai2002/his_v1/internal/domain"
)

func TestGender_Value(t *testing.T) {
	tests := []struct {
		name     string
		gender   domain.Gender
		expected string
	}{
		{"Male gender", domain.Male, "M"},
		{"Female gender", domain.Female, "F"},
		{"Other gender", domain.Other, "OTHER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.gender.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestGender_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected domain.Gender
	}{
		{"Scan M", "M", domain.Male},
		{"Scan F", "F", domain.Female},
		{"Scan OTHER", "OTHER", domain.Other},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g domain.Gender
			err := g.Scan(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, g)
		})
	}
}

func TestBloodGrp_Value(t *testing.T) {
	tests := []struct {
		name     string
		bloodGrp domain.BloodGrp
		expected string
	}{
		{"Blood type A", domain.A, "A"},
		{"Blood type B", domain.B, "B"},
		{"Blood type O", domain.O, "O"},
		{"Blood type AB", domain.AB, "AB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.bloodGrp.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestBloodGrp_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected domain.BloodGrp
	}{
		{"Scan A", "A", domain.A},
		{"Scan B", "B", domain.B},
		{"Scan O", "O", domain.O},
		{"Scan AB", "AB", domain.AB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b domain.BloodGrp
			err := b.Scan(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, b)
		})
	}
}

func TestPatientPartialUpdateRequest_ToMap(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	dateOfBirth := "1990-01-15"
	gender := "M"

	tests := []struct {
		name        string
		request     *domain.PatientPartialUpdateRequest
		expectedLen int
		expectError bool
	}{
		{
			name: "single field update",
			request: &domain.PatientPartialUpdateRequest{
				FirstNameEN: &firstName,
			},
			expectedLen: 1,
			expectError: false,
		},
		{
			name: "multiple fields update",
			request: &domain.PatientPartialUpdateRequest{
				FirstNameEN: &firstName,
				LastNameEN:  &lastName,
				Gender:      &gender,
			},
			expectedLen: 3,
			expectError: false,
		},
		{
			name: "update with date of birth",
			request: &domain.PatientPartialUpdateRequest{
				FirstNameEN: &firstName,
				DateOfBirth: &dateOfBirth,
			},
			expectedLen: 2,
			expectError: false,
		},
		{
			name:        "empty request",
			request:     &domain.PatientPartialUpdateRequest{},
			expectedLen: 0,
			expectError: false,
		},
		{
			name: "invalid date format",
			request: &domain.PatientPartialUpdateRequest{
				DateOfBirth: func() *string { s := "invalid-date"; return &s }(),
			},
			expectedLen: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.request.ToMap()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}

func TestPatientPartialUpdateRequest_ToMap_DateParsing(t *testing.T) {
	dateStr := "1990-05-20"
	request := &domain.PatientPartialUpdateRequest{
		DateOfBirth: &dateStr,
	}

	result, err := request.ToMap()
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	dob, ok := result["date_of_birth"].(time.Time)
	assert.True(t, ok)
	assert.Equal(t, 1990, dob.Year())
	assert.Equal(t, time.May, dob.Month())
	assert.Equal(t, 20, dob.Day())
}

func TestPatientPartialUpdateRequest_ToMap_AllFields(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	middleName := "Middle"
	firstNameTH := "จอห์น"
	lastNameTH := "โด"
	middleNameTH := "มิดเดิล"
	dateOfBirth := "1990-01-15"
	nickNameTH := "จอห์นนี่"
	nickNameEN := "Johnny"
	nationalID := "1234567890123"
	passportID := "AB123456"
	phoneNumber := "0812345678"
	email := "john@example.com"
	gender := "M"
	nationality := "Thai"
	bloodGrp := "A"

	request := &domain.PatientPartialUpdateRequest{
		FirstNameEN:  &firstName,
		LastNameEN:   &lastName,
		MiddleNameEN: &middleName,
		FirstNameTH:  &firstNameTH,
		LastNameTH:   &lastNameTH,
		MiddleNameTH: &middleNameTH,
		DateOfBirth:  &dateOfBirth,
		NickNameTH:   &nickNameTH,
		NickNameEN:   &nickNameEN,
		NationalID:   &nationalID,
		PassportID:   &passportID,
		PhoneNumber:  &phoneNumber,
		Email:        &email,
		Gender:       &gender,
		Nationality:  &nationality,
		BloodGrp:     &bloodGrp,
	}

	result, err := request.ToMap()
	assert.NoError(t, err)
	assert.Len(t, result, 16)

	assert.Equal(t, firstName, result["first_name_en"])
	assert.Equal(t, lastName, result["last_name_en"])
	assert.Equal(t, gender, result["gender"])
	assert.Equal(t, email, result["email"])
}

func TestPatient_Struct(t *testing.T) {
	patient := domain.Patient{
		FirstNameTH:  "สมชาย",
		LastNameTH:   "ใจดี",
		FirstNameEN:  "Somchai",
		LastNameEN:   "Jaidee",
		DateOfBirth:  time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		PatientHN:    "HOSP0001-00000001",
		NationalID:   "1234567890123",
		PhoneNumber:  "0812345678",
		Email:        "somchai@example.com",
		Gender:       domain.Male,
		Nationality:  "Thai",
		BloodGrp:     domain.A,
	}

	assert.Equal(t, "สมชาย", patient.FirstNameTH)
	assert.Equal(t, "Somchai", patient.FirstNameEN)
	assert.Equal(t, domain.Male, patient.Gender)
	assert.Equal(t, domain.A, patient.BloodGrp)
	assert.Equal(t, "HOSP0001-00000001", patient.PatientHN)
}
