package paypal

import (
	"os"
	"testing"
)

func TestIPN_VerifyIPN(t *testing.T) {
	ipn, err := NewPaypalIPN()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := ipn.VerifyIPN(os.Getenv("PAYPAL_IPN_BODY"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
