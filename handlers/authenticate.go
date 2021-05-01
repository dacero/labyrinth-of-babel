package handlers

import (
	"net/http"
	"log"
	"os"
	"fmt"

	"github.com/gorilla/sessions"
)

func Authenticate(store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "lob-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		secret := r.FormValue("secret")
		if secret != os.Getenv("LABYRINTH_SECRET") {
			session.Values["authenticated"] = false
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		
		session.Values["authenticated"] = true
		log.Print("It worked!")
		err = session.Save(r, w)
		if err != nil {
			log.Print("Internal Server Error when saving the session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(w, r, "/cell/entry", http.StatusFound)
	})
}
