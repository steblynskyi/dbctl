{{ define "dbo.source" }}
UPDATE 
    [dbo].[source]
SET 
    active=-2 
WHERE 
    sourceURL like '%.client.steblynskyi.com' 
	AND active=1  
{{ end }}
