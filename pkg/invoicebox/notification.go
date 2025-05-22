package invoicebox

import (
	"encoding/json"
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

	MetaData struct {
		TgChatID int64 `json:"chatId"`
	} `json:"metaData"`
}

func (h *WebhookHandler) HandleConfirmation(w http.ResponseWriter, r *http.Request) (paymentStatus string, chatId int64) {
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

	h.Printf("Получен вебхук от InvoiceBox:")
	h.Printf("%+v\n", notification)
	h.Printf("'\nResponse Body:\n", string(body))

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

	chatId = notification.MetaData.TgChatID

	h.Printf("TgID from API (notification): %d\n", chatId)
	h.Printf("Ожидали %.2f, пришло %.2f\n", task.Budget, notification.Amount)

	if notification.Status == "completed" && notification.Amount == task.Budget {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
		h.Printf("amount match \n")
		paymentStatus = "success"
	} else {
		http.Error(w, "amount mismatch or invalid status", http.StatusBadRequest)
		h.Printf("amount mismatch \n")
		paymentStatus = "fail"
	}

	if err != nil {
		return
	}

	return
}
