package password_hash_rest

import (
	"net/http"
	"runtime"
	"strconv"

	models "github.com/L1z1ng3r-sswe/all_knowledge-authentication/app/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHasher(password string) (string, *models.Response) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.Response{err, "Internal Server Error", err.Error(), http.StatusInternalServerError, getFileInfo("password.go")}
	}

	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, password string) *models.Response {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return &models.Response{err, "Bad Request", "Wrong password", http.StatusBadRequest, getFileInfo("password.go")}
	}

	return nil
}

func getFileInfo(fileName string) string {
	_, _, line, _ := runtime.Caller(1)
	return "internal/rest/utils/password/" + fileName + " line: " + strconv.Itoa(line)
}
