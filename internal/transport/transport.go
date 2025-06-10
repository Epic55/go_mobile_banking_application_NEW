package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/epic55/BankAppNew/internal/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	service services.ServiceInterface
}

func NewHandler(service services.ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", h.GetUserByID).Methods("GET")
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
