package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID uuid.UUID `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"`
	// Test comment
	CreatedAt time.Time `gorm:"type:timestamp without time zone;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;not null;default:now()"`
}
