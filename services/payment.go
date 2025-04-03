package services

import (
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type PaymentService struct{}

// PaymentItem represents an item with corresponding quantity in a payment
type PaymentItem struct {
	StripePriceId string `json:"stripe_price_id"`
	Quantity      uint64 `json:"quantity"`
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// Create a new payment details for order
func (ps *PaymentService) CreateNewPayment(orderId uint64, email string, paymentItemList []PaymentItem) (string, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range paymentItemList {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.StripePriceId),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	expiresAt := time.Now().Add(30 * time.Minute).Unix() // default 15 mins to checkout

	// Create a new checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"paynow",
		}),
		LineItems:     lineItems,
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String("http://localhost:3000/stripe/success"),
		CancelURL:     stripe.String("http://localhost:3000/stripe/cancel"),
		ExpiresAt:     stripe.Int64(expiresAt),
		CustomerEmail: stripe.String(email),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %v", err.Error())
	}

	return sess.URL, nil
}
