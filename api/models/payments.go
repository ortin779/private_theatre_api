package models

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
