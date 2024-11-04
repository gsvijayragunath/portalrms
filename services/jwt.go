package services

import (
	"errors"
	"example.com/RMS/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func GenerateToken(email string, userID uuid.UUID, userType string) (string, error) {
	authtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	})
	token, err := authtoken.SignedString([]byte(db.AuthKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(token string) (string, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte(db.AuthKey), nil
	})

	if err != nil {
		return "", "", errors.New("could not parse token")
	}
	if !parsedToken.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("userID not found in token claims")
	}

	userType, ok := claims["user_type"].(string)
	if !ok {
		return "", "", errors.New("userType not found in token claims")
	}
	return userID, userType, nil
}
