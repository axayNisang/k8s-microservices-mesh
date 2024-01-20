package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TokenPayload struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
}

var secretKey = []byte("super-secure-production-secret-key-32b")

func generateToken(userID, role string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadData, _ := json.Marshal(TokenPayload{
		UserID:    userID,
		Role:      role,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	payload := base64.RawURLEncoding.EncodeToString(payloadData)

	signature := computeSignature(header + "." + payload)
	return fmt.Sprintf("%s.%s.%s", header, payload, signature), nil
}

func validateToken(tokenString string) (*TokenPayload, bool) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, false
	}
	expectedSig := computeSignature(parts[0] + "." + parts[1])
	if parts[2] != expectedSig {
		return nil, false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, false
	}

	var payload TokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, false
	}

	if time.Now().Unix() > payload.ExpiresAt {
		return nil, false
	}
	return &payload, true
}

func computeSignature(message string) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP","service":"auth-service"}`))
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		token, _ := generateToken("usr_demo_101", "developer")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   86400,
		})
	})

	http.HandleFunc("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		payload, valid := validateToken(token)
		w.Header().Set("Content-Type", "application/json")
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"valid":false,"error":"Invalid or expired token"}`))
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   true,
			"user_id": payload.UserID,
			"role":    payload.Role,
		})
	})

	log.Println("[AuthService] Listening on port 8081 (HTTP) and gRPC stub ready...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
