package models

import "time"

type Url struct {
	Url string `json:"url"`
	CustomShort string `json:"short"`
	Expiry time.Duration `json:"expiry"`
}

// type 