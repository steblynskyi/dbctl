{{ define "stripe.PaymentMetadata" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.PaymentMetadata
{{ end }}