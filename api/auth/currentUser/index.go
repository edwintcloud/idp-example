package currentUser

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")))

// Handler is the primary handler for this route
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	gob.Register(map[string]interface{}{})

	// Get a session.
	session, err := store.Get(r, "currentUser")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if v, ok := session.Values["user"]; ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(v)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode("user unauthorized")

}
