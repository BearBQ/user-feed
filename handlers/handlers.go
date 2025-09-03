package handlers

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("работает ёпта")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Привет скуфам"))
}

func FuckHandler(w http.ResponseWriter, r *http.Request) {

}
