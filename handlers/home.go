package handlers

import (
	"net/http"
	"text/template"
)

// HomeHandler обрабатывает запрос на главную страницу
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
