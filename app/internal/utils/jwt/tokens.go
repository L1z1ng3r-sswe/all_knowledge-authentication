package tokens_grpc

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	models "github.com/L1z1ng3r-sswe/all_knowledge-authentication/app/internal/domain"
	"github.com/golang-jwt/jwt"
)

func CreateAccessToken(userId int64, exp time.Duration, secretKey string) (string, *models.Response) {
	claims := jwt.MapClaims{"sub": userId, "exp": time.Now().Add(exp).Unix()}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", &models.Response{err, err.Error(), "Internal Server Error", http.StatusInternalServerError, getFileInfo("tokens.go")}
	}
	return tokenString, nil
}

func CreateRefreshToken(userId int64, exp time.Duration, secretKey string) (string, *models.Response) {
	claims := jwt.MapClaims{"sub": userId, "exp": time.Now().Add(exp).Unix()}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", &models.Response{err, err.Error(), "Internal Server Error", http.StatusInternalServerError, getFileInfo("tokens.go")}
	}
	return tokenString, nil
}

func IsTokenValid(tokenString string, secretKey string) (int64, *models.Response) {

	tokenSlice := strings.Split(tokenString, " ")

	if len(tokenSlice) != 2 || tokenSlice[0] != "Bearer" {
		return 0, &models.Response{errors.New("Invalid authorization header format"), "Unauthorized", "Invalid authorization header format", http.StatusUnauthorized, getFileInfo("tokens.go")}
	}

	token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: " + token.Header["alg"].(string))
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Token is malformed", http.StatusUnauthorized, getFileInfo("tokens.go")}
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Token has expired", http.StatusUnauthorized, getFileInfo("tokens.go")}
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Token not yet valid", http.StatusUnauthorized, getFileInfo("tokens.go")}
			} else {
				return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Couldn't handle this token", http.StatusUnauthorized, getFileInfo("tokens.go")}
			}
		}
		return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Couldn't handle this token", http.StatusUnauthorized, getFileInfo("tokens.go")}
	}

	if !token.Valid {
		return 0, &models.Response{errors.New("Invalid token"), "Unauthorized", "Invalid token", http.StatusUnauthorized, getFileInfo("tokens.go")}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, &models.Response{errors.New("Invalid token claims"), "Unauthorized", "Invalid token claims", http.StatusUnauthorized, getFileInfo("tokens.go")}
	}

	userIdFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, &models.Response{errors.New("Invalid userId in token"), "Unauthorized", "Invalid userId in token", http.StatusUnauthorized, getFileInfo("tokens.go")}
	}

	userId := int64(userIdFloat)
	return userId, nil
}

func getFileInfo(fileName string) string {
	_, _, line, _ := runtime.Caller(1)
	return "internal/rest/utils/jwt/" + fileName + " line: " + strconv.Itoa(line)
}
