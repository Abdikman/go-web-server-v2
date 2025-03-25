package handlers

import (
	"context"
	"encoding/json"
	"go-web-server-v2/database"
	"go-web-server-v2/models"
	"net/http"
	"text/template"
)

// HomeHandler обрабатывает запрос на главную страницу
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Достаём имя файла из контекста
	filename := r.Context().Value("templateFile").(string)

	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Find(&users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func WithTemplateFile(filename string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "templateFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
