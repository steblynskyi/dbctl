{{ define "payments.DeviceTransaction" }}
USE steblynskyiPayments

TRUNCATE TABLE Payments.DeviceTransaction
{{ end }}