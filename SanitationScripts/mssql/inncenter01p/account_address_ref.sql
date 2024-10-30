{{ define "dbo.account_address_ref" }}
SET NOCOUNT ON;

DECLARE @min INT = 0,
        @max INT,
        @interval INT = 10000,
        @maxAcc INT,
        @qaEMail NVARCHAR(20) = N'QA-DG@steblynskyi.com';

SELECT @maxAcc = MAX(a.accountID)
  FROM dbo.account_address_ref a;

SET @max = @interval;

WHILE 1 = 1
BEGIN

    UPDATE aar
       SET aar.email = @qaEMail
      FROM dbo.account_address_ref aar
     WHERE aar.accountID BETWEEN @min AND @max
       AND aar.email <> @qaEMail
       AND aar.email IS NOT NULL;

    IF @max > @maxAcc
      RETURN;

    SET @min = @max + 1;
    SET @max += @interval; 

END;

{{ end }}