{{ define "stripe.PaymentIntent" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.PaymentIntent
{{ end }}