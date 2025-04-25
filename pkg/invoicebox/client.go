package invoicebox

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-telegram/bot"
	"github.com/labstack/echo/v4"
	"io"
	"lotBot/pkg/db"
	"lotBot/pkg/embedlog"
	"net/http"
	"time"
)

type Config struct {
	Database *pg.Options
	Server   struct {
		Host      string
		Port      int
		IsDevel   bool
		EnableVFS bool
	}
	Bot struct {
		Token string
	}
	Invoice struct {
		SecretKey  string
		MerchantID string
	}
}

type App struct {
	embedlog.Logger
	appName string
	cfg     Config
	db      db.DB
	dbc     *pg.DB
	echo    *echo.Echo
	b       *bot.Bot
}

// generateSignature is a helper function that generates an HMAC SHA-256 signature.
func generateSignature(secretKey string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

const url = "https://api.invoicebox.ru/v3/billing/api/order/order"

func (a *App) AskApi() {

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
		SuccessURL      string       `json:"successUrl"`
		FailURL         string       `json:"failUrl"`
		ReturnURL       string       `json:"returnUrl"`
		VatAmount       float64      `json:"vatAmount"`
		BasketItems     []BasketItem `json:"basketItems"`
	}

	secretKey := "QPu8HGhZ4iuOpgfVxdPEmV7ct3NCQozv"

	order := CreateOrderRequest{
		MerchantOrderID: "order-1234567890",
		Amount:          22,
		SuccessURL:      "https://merchant.ru/order/xxx?result=success",
		FailURL:         "https://merchant.ru/order/xxx?result=fail",
		ReturnURL:       "https://merchant.ru/order/xxx?result=return",
		VatAmount:       123,
		BasketItems: []BasketItem{
			{
				SKU:            "sku123",
				Name:           "qweqwe",
				Measure:        "шт.",
				MeasureCode:    "796",
				GrossWeight:    0,
				NetWeight:      0,
				Quantity:       3,
				Amount:         22,
				AmountWoVat:    123,
				TotalAmount:    1234,
				TotalVatAmount: 123,
				VatCode:        "RUS_VAT20",
				Type:           "service",
				PaymentType:    "full_prepayment",
			},
		},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}

	signature := generateSignature(secretKey, jsonData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MyApp 1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}
