package helpers

import "github.com/RakhmanovTimur/bookings/internal/config"

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app := a
}

