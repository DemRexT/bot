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
	SecretKey  string
	MerchantID string
}

type InvoiceClient struct {
	embedlog.Logger
	cfg Config
}

func NewInvoiceClient(logger embedlog.Logger, cfg Config) *InvoiceClient {
	return &InvoiceClient{Logger: logger, cfg: cfg}
}

const url = "https://api.invoicebox.ru/l3/billing/api/order/order"

func (ic *InvoiceClient) AskApi() error {
	type BasketItem struct {
		SKU            string  `json:"sku"`
		Name           string  `json:"name"`
		Measure        string  `json:"measure"`
		MeasureCode    string  `json:"measureCode"`
		GrossWeight    float64 `json:"grossWeight"`
		NetWeight      float64 `json:"netWeight"`
		Quantity       float64 `json:"quantity"`
		Amount         float64 `json:"amount"`
		AmountWoVat    float64 `json:"amountWoVat"`
		TotalAmount    float64 `json:"totalAmount"`
		TotalVatAmount float64 `json:"totalVatAmount"`
		VatCode        string  `json:"vatCode"`
		Type           string  `json:"type"`
		PaymentType    string  `json:"paymentType"`
	}

	type CreateOrderRequest struct {
		MerchantID      string       `json:"merchantId"`
		MerchantOrderID string       `json:"merchantOrderId"`
		Amount          float64      `json:"amount"`
		CurrencyID      string       `json:"currencyId"`
		SuccessURL      string       `json:"successUrl"`
		FailURL         string       `json:"failUrl"`
		ReturnURL       string       `json:"returnUrl"`
		VatAmount       float64      `json:"vatAmount"`
		BasketItems     []BasketItem `json:"basketItems"`
	}

	order := CreateOrderRequest{
		MerchantID:      "44844f1e-4228-4bd2-bd9c-73f90e3e06ed",
		MerchantOrderID: "test-order-123",
		Amount:          100.00,
		CurrencyID:      "RUB",
		SuccessURL:      "https://merchant.ru/success",
		FailURL:         "https://merchant.ru/fail",
		ReturnURL:       "https://merchant.ru/return",
		VatAmount:       16.67,
		BasketItems: []BasketItem{
			{
				SKU:            "sku123",
				Name:           "Test Product",
				Measure:        "pcs",
				MeasureCode:    "796",
				GrossWeight:    0,
				NetWeight:      0,
				Quantity:       1,
				Amount:         100.00,
				AmountWoVat:    83.33,
				TotalAmount:    100.00,
				TotalVatAmount: 16.67,
				VatCode:        "RUS_VAT20",
				Type:           "service",
				PaymentType:    "full_prepayment",
			},
		},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("JSON marshal error", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request creation error", err)
		return err
	}
	fmt.Println("Authorization SecretKey:", ic.cfg.SecretKey)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MyApp 1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ic.cfg.SecretKey)

	fmt.Println("Request Headers:")
	fmt.Printf("Content-Type: %s\n", req.Header.Get("Content-Type"))
	fmt.Printf("User-Agent: %s\n", req.Header.Get("User-Agent"))
	fmt.Printf("Accept: %s\n", req.Header.Get("Accept"))
	fmt.Printf("Authorization: %s\n", req.Header.Get("Authorization"))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request error", err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read error", err)
		return err
	}

	fmt.Printf("Response Status:", resp.Status)
	fmt.Printf("Response Body:", string(body))
	fmt.Println(req.Header)

	return nil
}
