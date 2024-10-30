{{ define "dbo.crsRoomType_roomClass_xref" }}
select
  *
from
  dbo.crsRoomType_roomClass_xref
where
  roomclassid in (
    select
      roomclassid
    from
      roomclass_published_v
    where
      propertyid = 157
  )
{{ end }}