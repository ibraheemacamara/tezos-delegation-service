package db

import "time"

type Delegations struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	Delegator string    `gorm:"not null"`
	Timestamp time.Time `gorm:"ype:timestamp;not null"`
	Block     int32     `gorm:"not null"`
	Amount    int64     `gorm:"not null"`
}
