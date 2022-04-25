package models

import "time"

type IpDetail struct {
	token          string
	Count          int
	ExpirationTime time.Time
}

type TokenDetail struct {
	Token          string    `db:"invite_token"`
	Count          int       `db:"-"`
	ExpirationTime time.Time `db:"archived_at"`
}
