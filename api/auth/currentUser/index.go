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
	gob.Register(map[string]interface{}{})

	// Get a session.
	session, err := store.Get(r, "currentUser")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if v, ok := session.Values["user"]; ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(v)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode("user unauthorized")

}
