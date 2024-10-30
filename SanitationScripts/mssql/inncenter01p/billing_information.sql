{{ define "dbo.billing_information" }}
-- update billing_information set email = 'QA-DG@steblynskyi.com'
DECLARE @r INT 
SET
  @r = 1;

WHILE @r > 0 BEGIN
update
  TOP (10000) billing_information
set
  email = 'QA-DG@steblynskyi.com'
where
  email <> 'QA-DG@steblynskyi.com'
SET
  @r = @@ROWCOUNT;

END
{{ end }}