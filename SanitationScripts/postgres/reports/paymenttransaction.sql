{{ define "public.paymenttransaction" }}
UPDATE public.paymenttransaction
   SET nameoncard = NULL
     , last4carddigits = NULL
     , expdate = NULL
     , creditcardrawdata = NULL;
{{ end }}