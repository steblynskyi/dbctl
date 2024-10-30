{{ define "dbo.CrsProperty_Property_ref" }}
-- To disable all CRS properties except the test property

UPDATE Crs
set
  crs.crsstatus = IIF(crs.propertyid <> 157, 0, crs.crsstatus), ---- To disable all CRS properties except the test property
  crs.crspropertyid = IIF(crs.propertyid =157 AND crs.crsid = 3 , 'TSINNR', crs.crspropertyid), --  -- To set the password for test property
  crs.password =  IIF(crs.propertyid =157 AND crs.crsid = 3 , '1234', crs.password) --  -- To set the password for test property
FROM dbo.CrsProperty_Property_ref crs

{{ end }}