package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/NLstn/clubs/azure/acs"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret")

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SendMagicLinkEmail(email, link string) {
	acs.SendMail([]acs.Recipient{{Address: email}}, "Magic Link", "Click the link to login: "+link, "<a href='"+link+"'>Click here to login</a>")
}

func GenerateJWT(userID string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)
	return tokenStr
}
