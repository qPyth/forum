package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, data interface{}, code int) error {
	const op = "sendJson"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("%s, json decode error: %w", op, err)
	}
	return nil
}

func htmlResponse(w http.ResponseWriter, file string, data any, code int) error {
	if code != 200 {
		w.WriteHeader(code)
	}
	op := "htmlResponse"

	templates, err := template.ParseGlob("./ui/html/templates/*.html")
	if err != nil {
		return fmt.Errorf("%s parseGrob error: %w", op, err)
	}

	if err := templates.ExecuteTemplate(w, file, data); err != nil {
		return fmt.Errorf("%s execute template error: %w", op, err)
	}
	return nil
}
