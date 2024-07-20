package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	SecretKey          string
	AccessTokenExpiry  int
	RefreshTokenExpiry int
}

func loadEnvConfig() (*TokenConfig, error) {
	envTokenExpMins := os.Getenv("JWT_ACC_TOKEN_EXP_MINS")
	tokenExp, err := strconv.Atoi(envTokenExpMins)
	if err != nil {
		return nil, err
	}

	envRefresTokenExpMins := os.Getenv("JWT_REFRESH_TOKEN_EXP_MINS")
	refresTokenExp, err := strconv.Atoi(envRefresTokenExpMins)
	if err != nil {
		return nil, err
	}
	tokenConfig := TokenConfig{
		SecretKey:          os.Getenv("JWT_SECRET_KEY"),
		AccessTokenExpiry:  tokenExp,
		RefreshTokenExpiry: refresTokenExp,
	}
	return &tokenConfig, nil
}

type CustomClaims struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId string, roles []string) (string, error) {
	tokenConfig, err := loadEnvConfig()
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	claims := CustomClaims{
		UserId: userId,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(tokenConfig.AccessTokenExpiry))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenConfig.SecretKey))
}

func GenerateRefreshToken(userId string, roles []string) (string, error) {
	tokenConfig, err := loadEnvConfig()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	claims := CustomClaims{
		UserId: userId,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(tokenConfig.RefreshTokenExpiry))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenConfig.SecretKey))
}

func ValidateToken(tokenString string) error {
	claims, err := GetClaims(tokenString)
	if err != nil {
		return err
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}
	return nil
}

func GetClaims(tokenString string) (CustomClaims, error) {
	tokenConfig, err := loadEnvConfig()
	if err != nil {
		return CustomClaims{}, fmt.Errorf("validate access token: %w", err)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenConfig.SecretKey), nil
	})

	if err != nil {
		return CustomClaims{}, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return CustomClaims{}, fmt.Errorf("error while decoding the claims")
	}
	return *claims, nil
}
