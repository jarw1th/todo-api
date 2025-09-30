package jwttoken

import (
	recovery "ToDoProject/safety"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenDuration = 24 * time.Hour
var refreshTokenDuration = 7 * 24 * time.Hour

func GenerateTokens(userID int) (accessToken string, refreshToken string, err error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(accessTokenDuration).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString(recovery.GetJwt())
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(refreshTokenDuration).Unix(),
		"type":    "refresh",
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString(recovery.GetJwt())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return recovery.GetJwt(), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if !parsedToken.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid user_id", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RefreshToken(refreshToken string) (newAccessToken string, err error) {
	parsed, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return recovery.GetJwt(), nil
	})
	if err != nil || !parsed.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", fmt.Errorf("invalid refresh token claims")
	}

	userID := int(claims["user_id"].(float64))
	newAccess, _, err := GenerateTokens(userID)
	return newAccess, err
}
