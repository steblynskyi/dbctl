{{ define "stripe.PaymentFeeRefund" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.PaymentFeeRefund
{{ end }}