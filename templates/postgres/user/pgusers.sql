{{ define "PgGetAllUsers" }}
SELECT usename AS role_name,
 CASE
  WHEN usesuper AND usecreatedb THEN
    CAST('superuser, create database' AS pg_catalog.text)
  WHEN usesuper THEN
    CAST('superuser' AS pg_catalog.text)
  WHEN usecreatedb THEN
    CAST('create database' AS pg_catalog.text)
  ELSE
    CAST('' AS pg_catalog.text)
 END role_attributes
FROM pg_catalog.pg_user
ORDER BY role_name desc;
{{ end }}

{{ define "getPgVersion" }}
SELECT version();
{{ end }}

{{ define "getPgRoles" }}
SELECT rolname,rolcanlogin FROM pg_roles;
{{ end }}

{{ define "getAllPgDbNames" }}
SELECT datname FROM pg_database;
{{ end }}

{{ define "getPgServerRoleMembers" }}
SELECT r1.rolname as "role",r.rolname as member
FROM pg_catalog.pg_roles r LEFT JOIN pg_catalog.pg_auth_members m
ON (m.member = r.oid)
LEFT JOIN pg_roles r1 ON (m.roleid=r1.oid)                                  
WHERE r.rolcanlogin
ORDER BY 1;
{{ end }}