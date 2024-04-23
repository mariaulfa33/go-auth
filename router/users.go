package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "Get All User")
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "Get All User")
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("id:%s", userID)))

}
