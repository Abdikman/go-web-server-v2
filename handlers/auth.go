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
	"gorm.io/gorm"
)

// Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// var user models.User
	// json.NewDecoder(r.Body).Decode(&user)
	if r.Method != "POST" {
		http.Error(w, "Только метод POST поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим данные формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли уже такой пользователь
	var existingUser models.User
	result_username := database.DB.Where("username = ?", username).First(&existingUser)

	if result_username.Error == nil {
		http.Error(w, "❌ Пользователь с таким именем уже существует", http.StatusConflict)
		return
	} else if result_username.Error != nil && result_username.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Ошибка проверки пользователя", http.StatusInternalServerError)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
	}

	user := models.User{Username: username, Password: string(hashedPassword)}

	// Сохраняем пользователя в БД
	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Ошибка сохранения в БД", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Успешная регистрация!"))
	// json.NewEncoder(w).Encode(user)
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
