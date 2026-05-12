package handler

import (
	"encoding/json"
	"net/http"

	"bank-service/internal/service"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(
	s *service.CardService,
) *CardHandler {
	return &CardHandler{service: s}
}

func (h *CardHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	var req struct {
		AccountID int64 `json:"account_id"`
	}

	err := json.NewDecoder(r.Body).
		Decode(&req)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	card, pan, cvv, err := h.service.Create(
		userID,
		req.AccountID,
	)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"card_id": card.ID,
		"pan":     pan,
		"cvv":     cvv,
	})
}

func (h *CardHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	cards, err := h.service.GetAll(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(cards)
}
