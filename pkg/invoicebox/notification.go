package invoicebox

import (
	"encoding/json"
	"fmt"
	"io"
	"lotBot/pkg/embedlog"
	"net/http"
	"strconv"

	"lotBot/pkg/db"
)

type WebhookHandler struct {
	DB db.DB
	embedlog.Logger
	repo db.LotbotRepo
}

func NewWebhookHandler(DB db.DB, logger embedlog.Logger) *WebhookHandler {
	return &WebhookHandler{
		DB: DB, Logger: logger,
		repo: db.NewLotbotRepo(DB),
	}
}

type InvoiceNotification struct {
	Type       string  `json:"type"`
	ID         string  `json:"merchantOrderId"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CurrencyID string  `json:"currencyId"`
}

func (h *WebhookHandler) HandleConfirmation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	var notification InvoiceNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	fmt.Printf("Получен вебхук от InvoiceBox:")
	fmt.Printf("%+v\n", notification)

	taskID, err := strconv.Atoi(notification.ID)
	if err != nil {
		http.Error(w, "invalid task ID format", http.StatusBadRequest)
		return
	}

	task, err := h.repo.TaskByID(r.Context(), taskID)
	if err != nil {
		h.Errorf("Task search error %v", err)
		return
	}

	fmt.Printf("Ожидали %.2f, пришло %.2f\n", task.Budget, notification.Amount)

	if notification.Status == "completed" && notification.Amount == task.Budget {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	} else {
		http.Error(w, "amount mismatch or invalid status", http.StatusBadRequest)
	}
}
