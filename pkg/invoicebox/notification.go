package invoicebox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"lotBot/pkg/db"

	"github.com/go-pg/pg/v10"
)

type WebhookHandlerDependencies struct {
	DB *pg.DB
}

type InvoiceNotification struct {
	Type       string  `json:"type"`
	ID         string  `json:"merchantOrderId"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CurrencyID string  `json:"currencyId"`
}

func (h *WebhookHandlerDependencies) WebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	task := &db.Task{}
	err = h.DB.Model(task).
		Where("taskId = ?", taskID).
		Select()
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
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
