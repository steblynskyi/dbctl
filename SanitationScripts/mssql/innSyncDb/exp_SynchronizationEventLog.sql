{{ define "expedia.SynchronizationEventLog" }}
USE innSyncDb

DELETE FROM expedia.SynchronizationEventLog
{{ end }}