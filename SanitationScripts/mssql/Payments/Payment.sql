{{ define "stripe.Payment" }}
USE steblynskyiPayments

DELETE FROM stripe.Payment
{{ end }}