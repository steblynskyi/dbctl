{{ define "stripe.AccountLog" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.AccountLog
{{ end }}