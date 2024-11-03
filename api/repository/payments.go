package repository

import (
	"database/sql"
)

type PaymentsRepository interface {
	Create(orderId, status string) error
	Update(orderId, signature, paymentId string) error
}

type paymentsRepository struct {
	db *sql.DB
}

func NewPaymentsRepository(db *sql.DB) PaymentsRepository {
	return &paymentsRepository{
		db: db,
	}
}

func (pr *paymentsRepository) Create(orderId, status string) error {
	_, err := pr.db.Exec(`INSERT INTO payments(razorpay_order_id, status ,razorpay_payment_id, razorpay_signature)
		VALUES ($1, $2, '' ,'')
	`, orderId, status)

	return err
}

func (pr *paymentsRepository) Update(orderId, signature, paymentId string) error {
	_, err := pr.db.Exec(`
        UPDATE payments
        SET razorpay_signature=$2,
            razorpay_payment_id=$3,
            status='success'
        WHERE razorpay_order_id = $1;
    `, orderId, signature, paymentId)

	return err
}
