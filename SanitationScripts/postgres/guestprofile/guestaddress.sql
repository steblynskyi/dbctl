{{ define "public.guestaddress" }}
UPDATE public.guestaddress
   SET email = 'QA-DG@steblynskyi.com',
       phonenumber = '1234567890',
       alternatephonenumber = '1234567890';
{{ end }}