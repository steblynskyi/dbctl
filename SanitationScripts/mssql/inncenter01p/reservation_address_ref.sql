{{ define "dbo.reservation_address_ref" }}
-- update reservation_address_ref set email = 'QA-DG@steblynskyi.com'
SET NOCOUNT ON;

DECLARE @min INT = -100,
        @max INT,
        @interval INT = 10000,
        @maxRes INT,
        @qaEMail NVARCHAR(20) = N'QA-DG@steblynskyi.com';

SELECT @maxRes = MAX(rar.reservationId)
  FROM dbo.reservation_address_ref rar;

SET @max = @interval;

WHILE 1 = 1
BEGIN

    UPDATE rar
       SET rar.email = @qaEMail
      FROM dbo.reservation_address_ref rar
     WHERE rar.reservationID BETWEEN @min AND @max
       AND rar.email <> @qaEMail
       AND rar.email IS NOT NULL;
    
    IF @max > @maxRes
      RETURN;

    SET @min = @max + 1;
    SET @max += @interval; 

END;

{{ end }}
