package types

import "time"

// Contact
type Contact struct {
	FullName  string
	Email     string
	Phone     string
	CreatedAt time.Time
}
