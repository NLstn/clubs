package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NLstn/clubs/azure/acs"
	"github.com/NLstn/clubs/database"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret")

type contextKey string

const UserIDKey contextKey = "userID"

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SendMagicLinkEmail(email, link string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	return acs.SendMail([]acs.Recipient{{Address: email}}, "Magic Link", "Click the link to login: "+link, "<a href='"+link+"'>Click here to login</a>")
}

func generateJWT(userID string, expiration time.Duration) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("cannot generate JWT with empty userID")
	}

	tokenID := GenerateToken()

	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     tokenID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenStr, nil
}

func GenerateAccessToken(userID string) (string, error) {
	return generateJWT(userID, 15*time.Minute)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateJWT(userID, 30*24*time.Hour)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Default().Println("Missing or invalid Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Default().Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Default().Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Default().Println("Token validation failed")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Default().Println("Could not parse claims as MapClaims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if claims["user_id"] == nil {
			log.Default().Println("user_id claim is missing")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Default().Printf("user_id is not a string: %T", claims["user_id"])
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Default().Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Default().Printf("Token parsing error: %v", err)
		return "", fmt.Errorf("invalid token")
	}

	if !token.Valid {
		log.Default().Println("Token validation failed")
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Println("Could not parse claims as MapClaims")
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Default().Printf("user_id is not a string: %T", claims["user_id"])
		return "", fmt.Errorf("invalid user ID")
	}

	var expiresAt time.Time
	err = database.Db.Raw(`SELECT expires_at FROM refresh_tokens WHERE user_id = ? AND token = ?`, userID, token.Raw).Scan(&expiresAt).Error
	if err != nil {
		return "", err
	}
	if expiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh token expired")
	}

	return userID, nil
}
