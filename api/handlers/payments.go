package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

type PaymentsHandler struct {
	logger          *zap.Logger
	paymentsService service.RazorpayService
}

func NewPaymentHandler(logger *zap.Logger, paymentsService service.RazorpayService) *PaymentsHandler {
	return &PaymentsHandler{
		logger:          logger,
		paymentsService: paymentsService,
	}
}

func (paymentsHandler *PaymentsHandler) VerifyPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var paymentBody models.PaymentVerificationBody

		err := json.NewDecoder(r.Body).Decode(&paymentBody)

		if err != nil {
			paymentsHandler.logger.Error("bad request", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusBadRequest, "unable to parse payment verification body")
			return
		}

		err = paymentsHandler.paymentsService.VerifyPayment(paymentBody)
		if err != nil {
			log.Println(err)
			if errors.Is(err, service.ErrPaymentSignatureFailure) {
				paymentsHandler.logger.Error("bad request", zap.String("error", err.Error()))
				RespondWithError(w, http.StatusBadRequest, "payment info is invalid")
				return
			}
			paymentsHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "successfully verified payment information"})
	}
}
