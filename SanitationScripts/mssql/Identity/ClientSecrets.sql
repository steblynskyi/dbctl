{{ define "dbo.ClientSecrets" }}

USE "Identity"

UPDATE dbo.ClientSecrets SET [Value] = ''
{{ end }}