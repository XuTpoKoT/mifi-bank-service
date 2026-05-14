package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-service/internal/service"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	service *service.CreditService
}

func NewCreditHandler(
	service *service.CreditService,
) *CreditHandler {
	return &CreditHandler{
		service: service,
	}
}

type createCreditRequest struct {
	AccountID  int64   `json:"account_id"`
	Principal  float64 `json:"principal"`
	TermMonths int     `json:"term_months"`
}

func (h *CreditHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	var req createCreditRequest

	err := json.NewDecoder(r.Body).
		Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.AccountID <= 0 {
		http.Error(
			w,
			"invalid account_id",
			http.StatusBadRequest,
		)
		return
	}

	if req.Principal <= 0 {
		http.Error(
			w,
			"principal must be > 0",
			http.StatusBadRequest,
		)
		return
	}

	if req.TermMonths <= 0 {
		http.Error(
			w,
			"term_months must be > 0",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Create(
		userID,
		req.AccountID,
		req.Principal,
		req.TermMonths,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "credit created",
		},
	)
}

func (h *CreditHandler) GetSchedule(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value("user_id").(int64)

	vars := mux.Vars(r)

	creditIDStr := vars["creditId"]

	creditID, err := strconv.ParseInt(
		creditIDStr,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid credit id",
			http.StatusBadRequest,
		)
		return
	}

	schedule, err := h.service.GetSchedule(
		userID,
		creditID,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		schedule,
	)
}
