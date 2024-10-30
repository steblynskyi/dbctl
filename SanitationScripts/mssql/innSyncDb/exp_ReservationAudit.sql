{{ define "expedia.ReservationAudit" }}
USE innSyncDb

DELETE FROM expedia.ReservationAudit
{{ end }}