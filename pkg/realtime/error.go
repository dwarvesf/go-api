package realtime

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrClientNotFound is returned when a client is not found.
	ErrClientNotFound = errors.New("client not found")

	// ErrDeviceNotFound is returned when a device is not found.
	ErrDeviceNotFound = errors.New("device not found")
)
