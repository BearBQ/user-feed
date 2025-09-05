package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"user-feed/db"
	"user-feed/models"
	"user-feed/redis"
)

type CustomHandler struct {
	db  db.DataBase
	rs  redis.RedisBase
	ctx context.Context
}

func NewCustomHandler(ctx context.Context, database db.DataBase, redis redis.RedisBase) *CustomHandler {
	return &CustomHandler{ctx: ctx, db: database, rs: redis}
}

func (h *CustomHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	var newUser models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.Users{
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	exist, err := h.db.CheckUserExist(newUser.Email)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if exist {
		http.Error(w, "Пользователь с таким Email существует", http.StatusConflict)
		return
	}

	err = h.db.CreateUser(&user)
	if err != nil {
		log.Printf("failed to create new user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Пользователь создан")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Пользователь успешно создан",
		"user_id": user.ID,
	})
}

// func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

// }

// func GetUserHandler(w http.ResponseWriter, r *http.Request) {
// 	userID := r.PathValue("userID")
// 	hellowMsg := fmt.Sprintf("Получил юзера %s", userID)
// 	log.Println(hellowMsg)
// 	w.Write([]byte(hellowMsg))

// }
