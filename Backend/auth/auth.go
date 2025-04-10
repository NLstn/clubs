package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/clubs/azure/acs"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtSecret = []byte("super-secret")

type contextKey string

const UserIDKey contextKey = "userID"

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SendMagicLinkEmail(email, link string) {
	acs.SendMail([]acs.Recipient{{Address: email}}, "Magic Link", "Click the link to login: "+link, "<a href='"+link+"'>Click here to login</a>")
}

func GenerateJWT(userID string) (string, error) {
	// Validate userID is not empty
	if userID == "" {
		return "", fmt.Errorf("cannot generate JWT with empty userID")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenStr, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Default().Println("Missing or invalid Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg
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

		// Extract user ID
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

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAuthorizedForClub(userId string, clubId string) bool {
	if userId == "" {
		return false
	}

	// Check if the user is the owner of the club
	var club models.Club
	result := database.Db.Where("id = ? AND owner_id = ?", clubId, userId).Find(&club)
	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		return false
	}
	if result.Error != nil {
		log.Default().Printf("Error checking club ownership: %v", result.Error)
		return false
	}
	// If the user is the owner, they are authorized

	return true
}
