{{ define "dbo.reservation" }}
-- update reservation set email = 'QA-DG@steblynskyi.com'
SET NOCOUNT ON;

DECLARE @min INT = -100,
        @max INT,
        @interval INT = 10000,
        @maxRes INT,
        @qaEMail NVARCHAR(20) = N'QA-DG@steblynskyi.com';

SELECT @maxRes = MAX(r.reservationId)
  FROM dbo.reservation r;

SET @max = @interval;

WHILE 1 = 1
BEGIN

    UPDATE res
       SET res.email = @qaEMail
      FROM dbo.reservation res
     WHERE res.reservationID BETWEEN @min AND @max
       AND res.email <> @qaEMail
       AND res.email IS NOT NULL;

    IF @max > @maxRes
      RETURN;

    SET @min = @max + 1;
    SET @max += @interval; 

END;

{{ end }}