package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtClaim struct {
	jwt.StandardClaims
	CustomerId int64 `json:"user_id"`
	Exp        int64 `json:"exp"`
}

type JWTToken struct {
	config *Config
}

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{config: config}
}
func (j *JWTToken) CreateToken(customer_id int64) (string, error) {
	//payload to encrypt
	claims := jwtClaim{
		CustomerId: customer_id,
		Exp:        time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.Signing_key))
	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Authentication token")
			//fmt.Errorf creates an error interface
		}
		return []byte(j.config.Signing_key), nil
	})
	if err != nil {
		return 0, fmt.Errorf("Invalid Authentication token")
	}
	claims, ok := token.Claims.(*jwtClaim)
	if !ok {
		return 0, fmt.Errorf("Invalid Authentication token")
	}
	if claims.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("token expired")
	}
	return claims.CustomerId, nil

}
