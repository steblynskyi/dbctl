{{ define "dbo.user_ref" }}
update
  dbo.user_ref
set
  email = 'QA-DG@steblynskyi.com'
where
  clientid <> 0
  and email not like '%steblynskyi.com%'
{{ end }}