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
	cfg Config
}

func NewInvoiceClient(logger embedlog.Logger, cfg Config) *InvoiceClient {
	return &InvoiceClient{Logger: logger, cfg: cfg}
}

const url = "https://api.invoicebox.ru/l3/billing/api/order/order"

func (ic *InvoiceClient) AskApi() (string, error) {
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

	type CreateOrderRequest struct {
		Description     string       `json:"description"`
		MerchantID      string       `json:"merchantId"`
		MerchantOrderID string       `json:"merchantOrderId"`
		Amount          float64      `json:"amount"`
		CurrencyID      string       `json:"currencyId"`
		VatAmount       float64      `json:"vatAmount"`
		BasketItems     []BasketItem `json:"basketItems"`
	}

	order := CreateOrderRequest{
		Description:     "Оплата услуг по оформлению бизнес-аккаунтов (Яндекс Бизнес, 2Гис)",
		MerchantID:      ic.cfg.MerchantID,
		MerchantOrderID: "1",
		Amount:          3000.00,
		CurrencyID:      "RUB",

		BasketItems: []BasketItem{
			{
				SKU:         "sku123",
				Name:        "Оформление бизнес-аккаунтов",
				Measure:     "шт.",
				Quantity:    1,
				Amount:      3000.00,
				VatCode:     "RUS_VAT20",
				Type:        "service",
				PaymentType: "full_prepayment",
			},
		},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("JSON marshal error", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request creation error", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MyApp 1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ic.cfg.SecretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request error", err)
		return "", err

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read error", err)
		return "", err
	}

	fmt.Printf("Response Status:", resp.Status)
	fmt.Printf("Response Body:", string(body))

	type CreateOrderResponse struct {
		Data struct {
			PaymentUrl string `json:"paymentUrl"`
		} `json:"data"`
	}

	var orderResp CreateOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	fmt.Println("paymentUrl from API:", orderResp.Data.PaymentUrl)

	if orderResp.Data.PaymentUrl == "" {
		return "", fmt.Errorf("paymentUrl not found in response")
	}

	return orderResp.Data.PaymentUrl, nil
}
