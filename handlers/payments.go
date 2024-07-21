package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ortin779/private_theatre_api/models"
)

func VerifyPayment(paymentService *models.RazorpayService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var paymentBody models.PaymentVerificationBody

		err := json.NewDecoder(r.Body).Decode(&paymentBody)

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusBadRequest, "unable to parse payment verification body")
			return
		}

		err = paymentService.VerifyPayment(paymentBody)
		if err != nil {
			log.Println(err)
			if errors.Is(err, models.ErrPaymentSingatureFailure) {
				RespondWithError(w, http.StatusBadRequest, "payment info is invalid")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "successfully verified payment information"})
	}
}
