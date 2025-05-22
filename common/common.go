package common

type PaymentStatusHandler interface {
	HandleStatus() (string, error)
}
