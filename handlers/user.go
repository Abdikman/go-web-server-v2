package handlers

import (
	"encoding/json"
	"go-web-server-v2/models"
	"net/http"
)

// Профиль пользователя (только для авторизованных)
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь можно будет получить данные о пользователе из токена
	user := models.User{Username: "Example"}
	json.NewEncoder(w).Encode(user)
}
