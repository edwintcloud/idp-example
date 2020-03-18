package callback

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
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
var store = sessions.NewCookieStore([]byte(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")))

// Handler is the primary handler for this route
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	gob.Register(map[string]interface{}{})

	// ensure state in callback matches state saved in cookie (to prevent CSRF)
	state, _ := r.Cookie("google-oauth-state")
	if r.FormValue("state") != state.Value {
		log.Printf("no dice \n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	// get user info from google user info api using state and code
	// passed in the callback
	userInfo, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	// set user in session
	session, err := store.Get(r, "currentUser")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{
		Path:     "/",       // the root of the app
		MaxAge:   86400 * 7, // 1 week
		HttpOnly: true,
	}
	session.Values["user"] = userInfo
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getUserInfo(state, code string) (map[string]interface{}, error) {
	var result map[string]interface{}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to convert code into token: %s", err.Error())
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("unable to get user info: %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("unable to decode resp.Body into map: %s", err.Error())
	}

	return result, nil
}
