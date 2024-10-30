{{ define "public.guestprofile_document" }}
UPDATE public.guestprofile_document
	   SET documentnumber = 'A123',
	       filepath = 'GuestProfile/0/0.png'
{{ end }}