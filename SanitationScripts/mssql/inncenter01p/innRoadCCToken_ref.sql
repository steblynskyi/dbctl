{{ define "dbo.steblynskyiCCToken_ref" }}
SET NOCOUNT ON;

DECLARE @min INT = 0,
        @max INT,
        @interval INT = 10000,
        @maxToken INT,
        @qaCCNum NVARCHAR(50) = N'aK0t2gRYWCdorS3aBFhYJwks04hYwi23';

SELECT @maxToken = MAX(t.steblynskyiCCTokenID)
  FROM dbo.steblynskyiCCToken_ref t;

SET @max = @interval;

WHILE 1 = 1
BEGIN

    UPDATE cct
       SET cct.CCNumber = @qaCCNum
      FROM dbo.steblynskyiCCToken_ref cct
     WHERE cct.steblynskyiCCTokenID BETWEEN @min AND @max
       AND cct.CCNumber <> @qaCCNum;

    IF @max > @maxToken
      RETURN;

    SET @min = @max + 1;
    SET @max += @interval; 

END;

{{ end }}