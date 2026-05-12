package handler

import (
	"bank-service/internal/service"
	"encoding/json"
	"net/http"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(
	s *service.AccountService,
) *AccountHandler {
	return &AccountHandler{s}
}

func (h *AccountHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	account, err := h.service.Create(userID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) TopUp(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	var req struct {
		AccountID int64   `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.service.TopUp(
		userID,
		req.AccountID,
		req.Amount,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (h *AccountHandler) Transfer(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	var req struct {
		FromID int64   `json:"from_account_id"`
		ToID   int64   `json:"to_account_id"`
		Amount float64 `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.service.Transfer(
		userID,
		req.FromID,
		req.ToID,
		req.Amount,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (h *AccountHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	accounts, err := h.service.GetAll(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}
