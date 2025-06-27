// Package auth ...
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/database"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type (
	// Repository ...
	Repository struct {
		Cfg *config.AppConfig
	}
)

// NewAuthRepository ...
func NewAuthRepository(cfg *config.AppConfig) (*Repository, error) {
	return &Repository{
		Cfg: cfg,
	}, nil
}

// CreateToken ...
func (u *Repository) CreateToken(data database.LoginDataWithRole, tokenType string) (string, error) {
	claims := jwt.MapClaims{}
	claims["mode"] = map[string]interface{}{
		"key":   tokenType,
		"value": true,
	}
	claims["user_id"] = data.UserID
	claims["role"] = data.UserWithRole.Role.Description
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // set token expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.Cfg.Secret))
}

func (u *Repository) VerifyToken(r *http.Request) (string, error) {
	tokenString := ""
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		tokenString = strings.Split(bearerToken, " ")[1]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(u.Cfg.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return pretty(claims), nil
	}
	return "", errors.New("token not valid")
}

func pretty(data interface{}) string {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Errorf("Failed to pretty claims %s", err)
		return ""
	}

	return string(b)
}
