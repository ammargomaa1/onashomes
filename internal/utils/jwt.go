package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/onas/ecommerce-api/config"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type EntityType string

const (
	EntityUser  EntityType = "user"
	EntityAdmin EntityType = "admin"
)

type Claims struct {
	EntityID   int64      `json:"entity_id"`
	EntityType EntityType `json:"entity_type"`
	RoleID     *int64     `json:"role_id,omitempty"`
	TokenType  TokenType  `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user or admin
func GenerateToken(entityID int64, entityType EntityType, roleID *int64, tokenType TokenType) (string, error) {
	cfg := config.AppConfig
	var expiryDuration time.Duration
	var secret string

	if tokenType == AccessToken {
		expiryDuration = time.Duration(cfg.JWT.Expiry) * time.Second
		secret = cfg.JWT.Secret
	} else {
		expiryDuration = time.Duration(cfg.JWT.RefreshExpiry) * time.Second
		secret = cfg.JWT.RefreshSecret
	}

	claims := Claims{
		EntityID:   entityID,
		EntityType: entityType,
		RoleID:     roleID,
		TokenType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	cfg := config.AppConfig
	var secret string

	if tokenType == AccessToken {
		secret = cfg.JWT.Secret
	} else {
		secret = cfg.JWT.RefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
