{{ define "public.batchreport" }}
UPDATE batchreport
   SET emaillist = '[]'
, istriggeremail = false;
{{ end }}