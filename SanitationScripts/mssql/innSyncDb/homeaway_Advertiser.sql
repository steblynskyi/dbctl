{{ define "homeaway.Advertiser" }}
USE innSyncDb

UPDATE [homeaway].[Advertiser] 
    SET [RentalAgreementUrl] = null, Active = 0
{{ end }}