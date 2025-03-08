package enum

import (
	"database/sql/driver"
	"fmt"
)

type PaymentType string

const (
	PaymentTypeUnpaid            PaymentType = "Unpaid"
	PaymentTypeQRPromptPay       PaymentType = "QRPromptPay"
	PaymentTypeCash              PaymentType = "Cash"
	PaymentTypeCreditOrDebitCard PaymentType = "CreditOrDebitCard"
)

func (r *PaymentType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*r = PaymentType(str)
	return nil
}

func (r *PaymentType) Value() (driver.Value, error) {
	return string(*r), nil
}
