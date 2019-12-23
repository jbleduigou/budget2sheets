package authentication

import (
	"fmt"
	"os"
)

const template string = `
{"installed":
  {"client_id":"%s",
    "project_id":"%s",
    "auth_uri":"https://accounts.google.com/o/oauth2/auth",
    "token_uri":"https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
    "client_secret":"%s",
    "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}

`

// GetCredentials will return json string for Google API credentials
func GetCredentials() string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	return fmt.Sprintf(template, clientID, projectID, clientSecret)
}
