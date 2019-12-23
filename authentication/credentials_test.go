package authentication

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expected string = `
{"installed":
  {"client_id":"8Pmwb8sLvzK72NTD.apps.googleusercontent.com",
    "project_id":"3i8N0VsGZQv0oic3",
    "auth_uri":"https://accounts.google.com/o/oauth2/auth",
    "token_uri":"https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
    "client_secret":"ke6hdeVUbHUJ6vji",
    "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}

`

func TestGetCredentials(t *testing.T) {
	os.Setenv("GOOGLE_CLIENT_ID", "8Pmwb8sLvzK72NTD.apps.googleusercontent.com")
	os.Setenv("GOOGLE_PROJECT_ID", "3i8N0VsGZQv0oic3")
	os.Setenv("GOOGLE_CLIENT_SECRET", "ke6hdeVUbHUJ6vji")

	output := GetCredentials()
	assert.Equal(t, expected, output)

	os.Unsetenv("GOOGLE_CLIENT_ID")
	os.Unsetenv("GOOGLE_PROJECT_ID")
	os.Unsetenv("GOOGLE_CLIENT_SECRET")
}
