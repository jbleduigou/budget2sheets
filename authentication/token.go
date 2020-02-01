package authentication

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// GetToken returns an oauth2 token for Google APIs
func GetToken() *oauth2.Token {
	tok := &oauth2.Token{}
	accessToken64, ok := os.LookupEnv("GOOGLE_ACCESS_TOKEN")
	if ok {
		decoded, err := base64.StdEncoding.DecodeString(accessToken64)
		if err != nil {
			fmt.Println(err)
		}
		tok.AccessToken = string(decoded)
	}
	tok.TokenType = "Bearer"
	tok.RefreshToken = os.Getenv("GOOGLE_REFRESH_TOKEN")
	tok.Expiry = time.Now()
	return tok
}
