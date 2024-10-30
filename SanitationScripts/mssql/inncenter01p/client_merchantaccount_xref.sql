{{ define "dbo.client_merchantaccount_xref" }}
-- To enable test mode on for merchant accounts so that the CC are not processed Live
UPDATE dbo.client_merchantAccount_xref
   SET testMode = 1,
       email = 'QA-DG@steblynskyi.com',
       status = IIF(gatewayID = 8 AND status <> -2, -2, status);
{{ end }}

/* 
Need to merge below query and combine into single query with above query

https://bitbucket.org/steblynskyi/databasesanitationscripts/src/1d001544581148d28405814ade20cae553b90581/SanitationScript-Inncenter.sql#lines-203:215

USE inncenter01p
GO
UPDATE 
     inncenter01p.dbo.client_merchantAccount_xref 
SET 
     STATUS=-2 
WHERE 
     gatewayID=8 
	 AND status<>-2 


*/
