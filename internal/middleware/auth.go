package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type AuthMiddleware struct {
	username string
	password string
	tokens   map[string]time.Time
	mu       sync.RWMutex
}

func NewAuthMiddleware(username, password string) *AuthMiddleware {
	m := &AuthMiddleware{
		username: username,
		password: password,
		tokens:   make(map[string]time.Time),
	}

	// Запускаем горутину для очистки устаревших токенов
	go m.cleanupExpiredTokens()

	return m
}

func (m *AuthMiddleware) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (m *AuthMiddleware) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for token, expiry := range m.tokens {
			if now.After(expiry) {
				delete(m.tokens, token)
			}
		}
		m.mu.Unlock()
	}
}

func (m *AuthMiddleware) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if credentials.Username != m.username || credentials.Password != m.password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := m.generateToken()
	m.mu.Lock()
	m.tokens[token] = time.Now().Add(24 * time.Hour)
	m.mu.Unlock()

	// ✅ Ставим cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true, если HTTPS
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
}
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		m.mu.RLock()
		expiry, exists := m.tokens[token]
		m.mu.RUnlock()

		if !exists || time.Now().After(expiry) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
