package config

import (
	"os"
)

// Configuration provides and interface for the configuration of the software
type Configuration interface {
	GetSpreadSheetID() string
	GetWriteRange() string
	GetGoogleClientID() string
	GetGoogleProjectID() string
	GetGoogleClientSecret() string
}

// NewConfiguration will provide an instance of a Configuration, implementation is not exposed
func NewConfiguration() Configuration {
	return &configuration{
		spreadSheetID:      os.Getenv("GOOGLE_SPREADSHEET_ID"),
		writeRange:         os.Getenv("GOOGLE_SPREADSHEET_RANGE"),
		googleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		googleProjectID:    os.Getenv("GOOGLE_PROJECT_ID"),
		googleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
}

type configuration struct {
	spreadSheetID      string
	writeRange         string
	googleClientID     string
	googleProjectID    string
	googleClientSecret string
}

func (c *configuration) GetSpreadSheetID() string {
	return c.spreadSheetID
}

func (c *configuration) GetWriteRange() string {
	return c.writeRange
}

func (c *configuration) GetGoogleClientID() string {
	return c.GetGoogleClientID()
}

func (c *configuration) GetGoogleProjectID() string {
	return c.GetGoogleProjectID()
}

func (c *configuration) GetGoogleClientSecret() string {
	return c.GetGoogleClientSecret()
}
