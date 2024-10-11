package model

import "time"

type User struct {
	LINEID  string    `firestore:"line_id"`
	Target  Target    `firestore:"target"`
	AddDate time.Time `firestore:"add_date"`
}

type PhoneNumber string

type Target struct {
	Nickname string      `firestore:"nickname"`
	Phone    PhoneNumber `firestore:"phone_number"`
	CallTime string      `firestore:"call_time"` // 12:00
	Confirm  bool        `firestore:"confirm"`
}
