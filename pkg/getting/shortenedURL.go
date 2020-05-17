package getting

import "time"

//ShortenedURL contains details for redirecting from hash key to full URL
type ShortenedURL struct {
	HashKey       string
	FullURL       string
	ExpirationUTC time.Time
}
