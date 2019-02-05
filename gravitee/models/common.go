package models

import "time"

type BaseModel struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TimestampModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type EmailTokenModel struct {
	BaseModel
	Reference   string `sql:"type:varchar(40):unique;not null"`
	EmailSet    bool   `sql:"index;not null"`
	EmailSentAt *time.Time
	ExpiresAt   time.Time `sql:"index;not null"`
}
