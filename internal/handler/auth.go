package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/abdul-hamid-achik/chessdrill/internal/service"
	"github.com/abdul-hamid-achik/chessdrill/templates/pages"
)

type AuthHandler struct {
	authService *service.AuthService
	maxAge      int
}

func NewAuthHandler(authService *service.AuthService, maxAge int) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		maxAge:      maxAge,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.Register("Invalid form data").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		pages.Register("All fields are required").Render(r.Context(), w)
		return
	}

	if len(password) < 6 {
		pages.Register("Password must be at least 6 characters").Render(r.Context(), w)
		return
	}

	_, token, err := h.authService.Register(r.Context(), email, username, password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			pages.Register("User already exists").Render(r.Context(), w)
			return
		}
		pages.Register("Registration failed").Render(r.Context(), w)
		return
	}

	h.setSessionCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.Login("Invalid form data").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		pages.Login("Email and password are required").Render(r.Context(), w)
		return
	}

	_, token, err := h.authService.Login(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			pages.Login("Invalid email or password").Render(r.Context(), w)
			return
		}
		pages.Login("Login failed").Render(r.Context(), w)
		return
	}

	h.setSessionCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.authService.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   h.maxAge,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(h.maxAge) * time.Second),
	})
}
