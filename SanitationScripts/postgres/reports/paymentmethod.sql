{{ define "public.paymentmethod" }}
UPDATE public.paymentmethod
   SET cclast4digits = NULL, expdate = NULL;
{{ end }}