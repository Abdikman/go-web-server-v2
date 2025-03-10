package handlers

import (
	"encoding/json"
	"go-web-server-v2/config"
	"go-web-server-v2/database"
	"go-web-server-v2/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// Хешируем пароль
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// Сохраняем пользователя в БД
	database.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Авторизация (JWT)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	database.DB.Where("username = ?", input.Username).First(&user)

	// Проверяем пароль
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := config.Config("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

}
