package token

import (
	"blog-web3/internal/configs"
	"blog-web3/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type Claims struct {
	model.User
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

func GenerateJWT(user *model.User) (string, error) {
	claims := Claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30000000 * time.Second).Unix(),
			Issuer:    "zou8944",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(configs.Conf.JWT.SignKey))
	return token, errors.WithStack(err)
}

func ParseJWT(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Conf.JWT.SignKey), nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tokenClaims.Claims.(*Claims), nil
}
