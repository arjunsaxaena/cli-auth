package cli

import "time"

type State struct {
	LoggedIn         bool
	UserID           int64
	Username         string
	SessionExpiresAt time.Time
}
