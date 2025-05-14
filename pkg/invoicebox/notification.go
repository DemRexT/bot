package invoicebox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InvoiceNotification struct {
	Type   string  `json:"type"`
	TaskID string  `json:"merchantOrderId"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
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

	fmt.Println("Получен вебхук от InvoiceBox:")
	fmt.Println(string(body))

	var notification InvoiceNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON: %v\n", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	fmt.Printf("Тип: %s\n", notification.Type)
	fmt.Printf("Статус: %s\n", notification.Status)
	fmt.Printf("Сумма: %s %s\n", notification.Amount, notification.Currency)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
