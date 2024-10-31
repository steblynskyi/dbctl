{{ define "stripe.Customer" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.Customer
{{ end }}