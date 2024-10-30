{{ define "expedia.SynchronizationEvent" }}
USE innSyncDb

TRUNCATE TABLE expedia.SynchronizationEvent
{{ end }}