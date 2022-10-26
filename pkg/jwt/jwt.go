package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/project5e/web3-blog/app/models"
	"github.com/project5e/web3-blog/config"
	"github.com/project5e/web3-blog/pkg/logger"
	"time"
)

type Claims struct {
	models.User
	jwt.StandardClaims
}

type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

func GenerateJWT(user *models.User) string {
	claims := Claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30000000 * time.Second).Unix(),
			Issuer:    "zou8944",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWT.SignKey))
	if err != nil {
		logger.Errorf("Generate jwt fail. %v", err)
		return ""
	}
	return token
}

func ParseJWT(token string) *Claims {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SignKey), nil
	})
	if err != nil {
		logger.Errorf("Parse jwt fail. %v", err)
		return nil
	}
	return tokenClaims.Claims.(*Claims)
}
