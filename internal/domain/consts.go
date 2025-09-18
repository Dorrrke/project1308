package domain

import "time"

// Timeouts.
const (
	ContextTimeout    = 6 * time.Second
	LeewayTimeout     = 60 * time.Second
	AccessTTL         = 15 * time.Minute
	RefreshTTL        = 24 * 7 * time.Hour
	CookieMaxAge      = 3600 * 24 * 7
	ReadHeaderTimeout = 15 * time.Second
)
