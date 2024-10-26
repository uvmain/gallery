package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

var (
	sessionToken = make(map[string]bool) // Track active sessions
	username     = os.Getenv("ADMIN_USER")
	password     = os.Getenv("ADMIN_PASSWORD")
)

// Helper function to generate a session token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("User logging in")
	if r.FormValue("username") == username && r.FormValue("password") == password {
		token := generateToken()
		sessionToken[token] = true
		http.SetCookie(w, &http.Cookie{
			Name:     "appSession",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		log.Println("Login successful")
		w.Write([]byte("Login successful"))
	} else {
		log.Println("Login unsuccessful, invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

// Middleware for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("appSession")
		if err != nil || !sessionToken[cookie.Value] {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Logout handler
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
