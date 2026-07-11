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

/*
type Base struct {
	ID uuid.UUID `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"`
	// Test comment
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
}
*/

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
	Users Users `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:SET NULL;"`
}

type SimpleTable struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name  string    `gorm:"type:varchar(255);not null"`
	SID   uint16    `gorm:"index;not null"`
	Email string    `gorm:"type:varchar(255);not null"`
	//UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Status    string    `gorm:"type:varchar(10);default:active"`
	NameT     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
	//Users     Users     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:SET NULL;"`
}

type Clicks struct {
	Base
	Type   string
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Users  Users     `gorm:"foreignKey:UserID;referenceName:clicks_users_user_id_id;constraint:OnDelete:CASCADE;"`
}

type Cars struct {
	Base
	Name            string `gorm:"type:varchar(255);not null"`
	Color           string `gorm:"type:varchar(255);default:grey"`
	EnginePowerHp   uint16 `gorm:"null"`
	Shifts          uint8  `gorm:"type:smallint;null"`
	ElectricPowerKw uint16 `gorm:"null"`
}

type Commands struct {
	Base
	Name  string    `gorm:"type:varchar(255);not null"`
	Cid   uuid.UUID `gorm:"type:uuid;index;default:uuidv4()"`
	CarId uuid.UUID `gorm:"type:uuid;not null"`
	Owner string    `gorm:"type:text;null"`
	//Description string    `gorm:"type:text;null"`
	BudgetM float32 `gorm:"null"`
	Cars    Cars    `gorm:"foreignKey:CarId;referenceName:commands_cars_car_id_id;constraint:OnDelete:CASCADE;"`
}
