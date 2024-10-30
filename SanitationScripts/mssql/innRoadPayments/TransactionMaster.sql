{{ define "payments.TransactionMaster" }}
USE steblynskyiPayments

TRUNCATE TABLE Payments.TransactionMaster
{{ end }}