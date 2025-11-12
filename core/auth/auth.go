package auth

import (
	"crypto/rand"
	"encoding/base64"
	"gallery/core/config"
	"log"
	"net/http"
)

var sessionToken = make(map[string]bool)

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	passedUsername := r.FormValue("username")
	passedPassword := r.FormValue("password")
	log.Printf("User attempting to log in: %s", passedUsername)

	if passedUsername == config.AdminUser && passedPassword == config.AdminPassword {
		token := generateToken()
		sessionToken[token] = true
		http.SetCookie(w, &http.Cookie{
			Name:     "appSession",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
			MaxAge:   604800,
		})
		log.Printf("Login successful for user: %s", passedUsername)
		w.Write([]byte("Login successful"))
	} else {
		log.Printf("Login unsuccessful for user: %s", passedUsername)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("appSession")
		if err != nil || !sessionToken[cookie.Value] {
			log.Printf("Unauthorized access attempt for user: %s", r.FormValue("username"))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("User logging out")
	cookie, err := r.Cookie("appSession")
	if err == nil {
		delete(sessionToken, cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "appSession",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Logged out"))
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("appSession")
	if err != nil || !sessionToken[cookie.Value] {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
