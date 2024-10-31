{{ define "stripe.ChargebackToPaymentMap" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.ChargebackToPaymentMap
{{ end }}