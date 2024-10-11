package model

import "time"

type User struct {
	LINEID  string    `json:"line_id"`
	Target  Target    `json:"target"`
	AddDate time.Time `json:"add_date"`
}

type PhoneNumber string

type Target struct {
	Nickname string      `json:"nickname"`
	Phone    PhoneNumber `json:"phone_number"`
	CallTime string      `json:"call_time"` // 12:00
}
