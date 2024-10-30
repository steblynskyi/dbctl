{{ define "stripe.PaymentHistory" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.PaymentHistory
{{ end }}