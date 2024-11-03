package models

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

type OrderAddon struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderAddonDetails struct {
	Addon
	Quantity int `json:"quantity"`
}

type PaymentStatus string

var (
	Success PaymentStatus = "success"
	Failure PaymentStatus = "failure"
	Pending PaymentStatus = "pending"
)

type OrderParams struct {
	CustomerName  string       `json:"customer_name"`
	CustomerEmail string       `json:"customer_email"`
	PhoneNumber   string       `json:"phone_number"`
	TheatreId     string       `json:"theatre_id"`
	SlotId        string       `json:"slot_id"`
	NoOfPersons   int          `json:"no_of_persons"`
	TotalPrice    int          `json:"total_price"`
	OrderDate     time.Time    `json:"order_date"`
	Addons        []OrderAddon `json:"addons"`
}

func (op *OrderParams) Validate() map[string]string {
	errs := make(map[string]string)

	if len(op.PhoneNumber) == 0 {
		errs["phone_number"] = "phone number can not be empty"
	}
	if len(op.CustomerName) == 0 {
		errs["customer_name"] = "phone number can not be empty"
	}
	if !isEmailValid(op.CustomerEmail) {
		errs["customer_email"] = "invalid email address"
	}
	if _, err := uuid.Parse(op.SlotId); err != nil {
		errs["slot_id"] = "slot id must be a valid uuid"
	}
	if op.TotalPrice <= 0 {
		errs["total_price"] = "order value must be greater than zero"
	}
	for _, addon := range op.Addons {
		if _, err := uuid.Parse(addon.ID); err != nil {
			errs["addons"] = "addon id must be a valid uuid"
			break
		}
		if addon.Quantity <= 0 {
			errs["addons"] = "addon quantity should not be less than 1"
			break
		}
	}
	return errs
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type OrderDetails struct {
	ID             string              `json:"id"`
	CustomerName   string              `json:"customer_name"`
	CustomerEmail  string              `json:"customer_email"`
	PhoneNumber    string              `json:"phone_number"`
	NoOfPersons    int                 `json:"no_of_persons"`
	TotalPrice     int                 `json:"total_price"`
	OrderDate      time.Time           `json:"order_date"`
	Theatre        Theatre             `json:"theatre"`
	Slot           Slot                `json:"slot"`
	Addons         []OrderAddonDetails `json:"addons"`
	OrderedAt      time.Time           `json:"ordered_at"`
	PaymentDetails OrderPayment        `json:"payment_details"`
}

type Order struct {
	ID              string       `json:"id"`
	CustomerName    string       `json:"customer_name"`
	CustomerEmail   string       `json:"customer_email"`
	PhoneNumber     string       `json:"phone_number"`
	TheatreId       string       `json:"theatre_id"`
	SlotId          string       `json:"slot_id"`
	Addons          []OrderAddon `json:"addons"`
	NoOfPersons     int          `json:"no_of_persons"`
	TotalPrice      int          `json:"total_price"`
	OrderDate       time.Time    `json:"order_date"`
	OrderedAt       time.Time    `json:"ordered_at"`
	RazorpayOrderId string       `json:"razorpay_order_id"`
}
