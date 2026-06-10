package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Pending  Status = "pending"
)

type Base struct {
	ID uuid.UUID `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"`
	// Test comment
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
}

type Users struct {
	Base
	Name   string `gorm:"not null"`       // full name field
	Email  string `gorm:"not null"`       // User email field
	Status string `gorm:"default:active"` // user statuses
}

type AuthUsers struct {
	Base
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	Password string    `gorm:"type:text;not null"`
	// Constraints
	Users Users `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type SimpleTable struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	SID       uint16    `gorm:"index;not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Status    string    `gorm:"type:varchar(10);default:active"`
	NameT     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
}
