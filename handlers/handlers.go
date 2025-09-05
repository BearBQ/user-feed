package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (h *CustomHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	var data models.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	var dataToBd = models.Post{
		Content: data.Content,
		UserID:  data.UserID,
	}
	err = h.db.CreatePost(&dataToBd)
	if err != nil {
		log.Printf("failed to create new post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = h.rs.ClearData(h.ctx, data.UserID)
	if err != nil {
		log.Printf("redis delete error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "new post create",
		"userID":  data.UserID,
		"Post":    data.Content,
	})

}

func (h *CustomHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	id, err := strconv.ParseUint(userID, 10, 64)
	idUint := uint(id)
	if err != nil {
		log.Printf("failed to get user id")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	data, err := h.rs.GetDataFromRedis(h.ctx, idUint)
	if err != nil {
		log.Printf("redis delete error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

		return
	}
	dataFromDB, err := h.db.GetPosts(idUint)
	if err != nil {
		log.Printf("postgres get post error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var dataToRedis []byte
	dataToRedis, err = json.Marshal(dataFromDB)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = h.rs.PushData(h.ctx, dataFromDB.ID, dataToRedis)
	if err != nil {
		log.Printf("push to redis error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataToRedis)
}
