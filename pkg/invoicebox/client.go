package invoicebox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lotBot/pkg/embedlog"
	"net/http"
	"time"
)

type Config struct {
	SecretKey  string `toml:"SecretKey"`
	MerchantID string `toml:"MerchantID"`
}

type InvoiceClient struct {
	embedlog.Logger
	cfg      Config
	TgChatID int64
}

func NewInvoiceClient(logger embedlog.Logger, cfg Config) *InvoiceClient {
	return &InvoiceClient{Logger: logger, cfg: cfg}
}

const url = "https://api.invoicebox.ru/l3/billing/api/order/order"

func (ic *InvoiceClient) AskApi(ChatID int64, taskId string, description string, budget float64, name string, StudentChat int64, YougileId string) (string, error) {
	type BasketItem struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Measure     string  `json:"measure"`
		Quantity    float64 `json:"quantity"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		VatCode     string  `json:"vatCode"`
		PaymentType string  `json:"paymentType"`
	}

	type MetaData struct {
		TgChatID      int64  `json:"TgChatID"`
		StudentTgId   int64  `json:"StudentTgId"`
		YougileIdTask string `json:"YougileIdTask"`
	}

	type CreateOrderRequest struct {
		Description     string       `json:"description"`
		MerchantID      string       `json:"merchantId"`
		MerchantOrderID string       `json:"merchantOrderId"`
		Amount          float64      `json:"amount"`
		CurrencyID      string       `json:"currencyId"`
		VatAmount       float64      `json:"vatAmount"`
		BasketItems     []BasketItem `json:"basketItems"`
		Metadata        MetaData     `json:"metaData"`
	}

	order := CreateOrderRequest{
		Description:     description,
		MerchantID:      ic.cfg.MerchantID,
		MerchantOrderID: taskId,
		Amount:          budget,
		CurrencyID:      "RUB",

		BasketItems: []BasketItem{
			{
				SKU:         taskId,
				Name:        name,
				Measure:     "шт.",
				Quantity:    1,
				Amount:      budget,
				VatCode:     "RUS_VAT20",
				Type:        "service",
				PaymentType: "full_prepayment",
			},
		},
		Metadata: MetaData{
			TgChatID:      ChatID,
			StudentTgId:   StudentChat,
			YougileIdTask: YougileId,
		},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		ic.Printf("JSON marshal error", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		ic.Printf("Request creation error", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MyApp 1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ic.cfg.SecretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ic.Printf("Request error", err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ic.Printf("Read error", err)
		return "", err
	}

	ic.Printf("Response Status:\n", resp.Status)
	ic.Printf("'\nResponse Body:\n", string(body))

	type CreateOrderResponse struct {
		Data struct {
			PaymentUrl string `json:"paymentUrl"`
			MetaData   struct {
				TgChatID int64 `json:"chatId"`
			} `json:"metaData"`
		} `json:"data"`
	}

	var orderResp CreateOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	ic.TgChatID = orderResp.Data.MetaData.TgChatID
	fmt.Println("Meta from API (client):", ic.TgChatID)
	fmt.Println("paymentUrl from API:\n", orderResp.Data.PaymentUrl)

	if orderResp.Data.PaymentUrl == "" {
		return "", fmt.Errorf("paymentUrl not found in response")
	}

	return orderResp.Data.PaymentUrl, nil
}
