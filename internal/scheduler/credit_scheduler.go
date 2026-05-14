package scheduler

import (
	"bank-service/internal/service"
	"time"
)

func StartCreditScheduler(
	creditService *service.CreditService,
) {
	go func() {
		for {
			creditService.ProcessDuePayments()

			time.Sleep(
				12 * time.Hour,
			)
		}
	}()
}
