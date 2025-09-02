package handlers

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("работает ёпта")
}

func FuckHandler(w http.ResponseWriter, r *http.Request) {

}
