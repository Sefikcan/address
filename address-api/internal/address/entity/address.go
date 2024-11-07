package entity

import "time"

type Address struct {
	Id          int       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	FullAddress string    `json:"full_address"`
	UserId      string    `json:"user_id"`
}
