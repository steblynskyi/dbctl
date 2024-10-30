{{ define "dbo.systables" }}

USE steblynskyiEventsDb

DECLARE @SQL NVARCHAR(max)=''
SELECT 
    @Sql=@sql+ 'TRUNCATE TABLE  ['+ schema_name(schema_id)+'].['+name+']; ' 
FROM sys.tables 
WHERE name NOT IN ('Streams')
EXEC (@sql)

{{ end }}