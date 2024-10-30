--Create SQL-login sample query
-- CREATE LOGIN appInncenter 
--     WITH PASSWORD ='TestPassword9%',
--     DEFAULT_DATABASE = inncenter01p,
--     CHECK_EXPIRATION = ON,
--     CHECK_POLICY = ON

{{ define "createUser"}}
CREATE LOGIN {{ .Username }} 
    WITH PASSWORD ='{{ .Password }}',
    DEFAULT_DATABASE = {{ .DefaultDb }},
    CHECK_EXPIRATION = ON,
    CHECK_POLICY = ON
{{ end }}

{{ define "getAllUsers" }}
SELECT name,principal_id,create_date,default_database_name FROM sys.sql_logins;
{{ end }}

{{ define "getAllDbNames" }}
SELECT name, database_id, create_date FROM sys.databases;
{{ end }}