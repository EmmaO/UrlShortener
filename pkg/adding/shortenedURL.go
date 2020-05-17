package adding

import "time"

//ShortenedURL contains details for redirecting from hash key to full URL
type ShortenedURL struct {
	HashKey       string
	FullURL       string
	ExpirationUTC time.Time
}

//ShortenedURLRequest is a truncated ShortenedURL which exposes fields that can be set by clients
type ShortenedURLRequest struct {
	HashKey string
	FullURL string
}
