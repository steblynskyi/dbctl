{{ define "homeaway.Advertiser" }}
USE innSyncDb

UPDATE [homeaway].[Listing] 
    SET ApprovalStatusId=1 
   WHERE ApprovalStatusId=2
{{ end }}