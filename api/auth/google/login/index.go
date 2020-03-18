package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// Handler is the primary handler for this route
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// generate random state
	state := generateOauthState()
	// create cookie with state
	cookie := createOauthCookie(state)
	// set cookie
	http.SetCookie(w, &cookie)
	// get google login url
	loginURL := googleOauthConfig.AuthCodeURL(state)
	// redirect to loginURL - google sign-in page
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func generateOauthState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func createOauthCookie(state string) http.Cookie {
	exp := time.Now().Add(20 * time.Minute) // 20 mins

	return http.Cookie{Name: "google-oauth-state", Value: state, Expires: exp}
}
