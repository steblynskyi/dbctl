{{ define "bookingdotcom.SynchronizationEventLog" }}
USE innSyncDb

DELETE FROM bookingdotcom.SynchronizationEventLog
{{ end }}