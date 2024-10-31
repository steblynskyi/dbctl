{{ define "stripe.ChargebackBalanceTransfer" }}
USE steblynskyiPayments

TRUNCATE TABLE stripe.ChargebackBalanceTransfer
{{ end }}