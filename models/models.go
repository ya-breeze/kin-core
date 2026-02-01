package models

import (
	"time"

	"gorm.io/gorm"
)

// Family represents a tenant group in the system
type Family struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
}

// User represents an individual within a Family
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"unique;not null" json:"username"`
	PasswordHash string         `json:"-"`
	FamilyID     uint           `gorm:"index;not null" json:"family_id"`
	Family       Family         `gorm:"foreignKey:FamilyID" json:"family,omitempty"`
}

// TenantModel is a base struct that should be embedded in all models
// that require data isolation at the family level.
type TenantModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyID  uint           `gorm:"index;not null" json:"family_id"`
}
