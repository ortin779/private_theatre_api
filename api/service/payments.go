package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
	"github.com/razorpay/razorpay-go"
)

type RazorpayService struct {
	paymentRepo repository.PaymentsRepository
	config      models.RazorpayConfig
	client      *razorpay.Client
}

var (
	ErrPaymentSignatureFailure = errors.New("invalid payment signature")
)

func NewRazorpayService(paymentRepo repository.PaymentsRepository, paymentConfig models.RazorpayConfig) RazorpayService {
	return RazorpayService{
		paymentRepo: paymentRepo,
		config:      paymentConfig,
		client:      razorpay.NewClient(paymentConfig.Key, paymentConfig.Secret),
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

	err = paymentService.paymentRepo.Create(paymentOrderId, string(models.Pending))

	if err != nil {
		return "", fmt.Errorf("create payment order: %w", err)
	}
	return paymentOrderId, nil
}

func (paymentService *RazorpayService) VerifyPayment(verificationBody models.PaymentVerificationBody) error {
	isValidSignature := verifySignature(verificationBody.RazorpayOrderId, verificationBody.RazorpayPaymentId, verificationBody.RazorpaySignature, paymentService.config.Secret)
	if !isValidSignature {
		return ErrPaymentSignatureFailure
	}

	err := paymentService.paymentRepo.Update(verificationBody.RazorpayOrderId, verificationBody.RazorpaySignature, verificationBody.RazorpayPaymentId)

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
