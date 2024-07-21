package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

type RazorpayConfig struct {
	Key    string
	Secret string
}

type PaymentVerificationBody struct {
	RazorpayOrderId   string `json:"razorpay_order_id"`
	RazorpayPaymentId string `json:"razorpay_payment_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}

func (paymentBody PaymentVerificationBody) Validate() map[string]string {
	errs := make(map[string]string)

	if paymentBody.RazorpayOrderId == "" {
		errs["razorpay_order_id"] = "razorpay_order_id can not be empty"
	}
	if paymentBody.RazorpayPaymentId == "" {
		errs["razorpay_payment_id"] = "razorpay_payment_id can not be empty"
	}
	if paymentBody.RazorpaySignature == "" {
		errs["razorpay_signature"] = "razorpay_signature can not be empty"
	}
	return errs
}

type OrderPayment struct {
	PaymentVerificationBody
	Status PaymentStatus `json:"status"`
}

type RazorpayService struct {
	db     *sql.DB
	config RazorpayConfig
	client *razorpay.Client
}

var (
	ErrPaymentSingatureFailure = errors.New("invalid payment signature")
)

func NewRazorpayService(db *sql.DB, paymentConfig RazorpayConfig) *RazorpayService {
	return &RazorpayService{
		db:     db,
		config: paymentConfig,
		client: razorpay.NewClient(paymentConfig.Key, paymentConfig.Secret),
	}
}

func (paymentService *RazorpayService) CreateOrder(amount int) (string, error) {
	razorpayData := map[string]any{
		"amount":          amount,
		"currency":        "INR",
		"partial_payment": false,
	}

	razorpayOrder, err := paymentService.client.Order.Create(razorpayData, nil)
	if err != nil {
		return "", fmt.Errorf("create payment order: %w", err)
	}

	paymentOrderId := (razorpayOrder["id"]).(string)

	_, err = paymentService.db.Exec(`INSERT INTO payments(razorpay_order_id, status ,razorpay_payment_id, razorpay_signature)
		VALUES ($1, $2, '' ,'')
	`, razorpayOrder["id"], Pending)

	if err != nil {
		return "", fmt.Errorf("create payment order: %w", err)
	}
	return paymentOrderId, nil
}

func (paymentService *RazorpayService) VerifyPayment(verificationBody PaymentVerificationBody) error {
	isValidSignature := verifySignature(verificationBody.RazorpayOrderId, verificationBody.RazorpayPaymentId, verificationBody.RazorpaySignature, paymentService.config.Secret)
	if !isValidSignature {
		return ErrPaymentSingatureFailure
	}

	_, err := paymentService.db.Exec(`
        UPDATE payments
        SET razorpay_signature=$2,
            razorpay_payment_id=$3,
            status='success'
        WHERE razorpay_order_id = $1;
    `, verificationBody.RazorpayOrderId, verificationBody.RazorpaySignature, verificationBody.RazorpayPaymentId)

	if err != nil {
		return fmt.Errorf("verify order payment: %w", err)
	}
	return nil
}

func verifySignature(orderId, paymentId, signature, secret string) bool {
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))

	if err != nil {
		return false
	}
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(generatedSignature), []byte(signature)) == 1
}
