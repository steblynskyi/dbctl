{{ define "dbo.property_ref" }}
-- To disable sending PBX messages
USE inncenter01p

UPDATE dbo.property_ref
   SET sendCCMail = 'QA-DG@steblynskyi.com',
       sendScheduleEmailsFrom = 'QA-DG@steblynskyi.com',
       emailFrmPropSelVal = 'QA-DG@steblynskyi.com', ---- To change the status of the pending items as completed so that these messages will not be processed by the steblynskyiserviceengine
       ipAddress = NULL,
       houseAccountID = NULL;
{{ end }}