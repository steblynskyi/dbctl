{{ define "createRole" }}
CREATE ROLE .RoleName 
{{ end }}


{{ define "getDBRoles" }}
USE {{ .DbName }}
SELECT name,principal_id,create_date FROM sys.database_principals where type='R'
{{ end }}


{{ define "getServerRoles" }}
SELECT name,principal_id FROM sys.server_principals where type='R';
{{ end }}

{{ define "getServerRoleMembers" }}
SELECT	roles.principal_id							AS RolePrincipalID
	,	roles.name									AS RolePrincipalName
	,	server_role_members.member_principal_id		AS MemberPrincipalID
	,	members.name								AS MemberPrincipalName
FROM sys.server_role_members AS server_role_members
INNER JOIN sys.server_principals AS roles
    ON server_role_members.role_principal_id = roles.principal_id
INNER JOIN sys.server_principals AS members 
    ON server_role_members.member_principal_id = members.principal_id
where roles.name='{{ .RoleName }}';
{{ end }}


{{ define "getDBRoleMembers" }}
USE {{ .DbName }}
SELECT    roles.principal_id                            AS RolePrincipalID
    ,    roles.name                                    AS RolePrincipalName
    ,    database_role_members.member_principal_id    AS MemberPrincipalID
    ,    members.name                                AS MemberPrincipalName
FROM sys.database_role_members AS database_role_members  
JOIN sys.database_principals AS roles  
    ON database_role_members.role_principal_id = roles.principal_id  
JOIN sys.database_principals AS members  
    ON database_role_members.member_principal_id = members.principal_id
where roles.name='{{ .RoleName }}';
{{ end }}


{{ define "getServerRolePermissions" }}
SELECT pr.principal_id, pr.name, pr.type_desc,   
    pe.state_desc, pe.permission_name   
FROM sys.server_principals AS pr 
JOIN sys.server_permissions AS pe   
    ON pe.grantee_principal_id = pr.principal_id
where pr.name='{{ .RoleName }}';
{{ end }}


{{ define "getDBRolePermissions" }}
USE {{ .DbName }}
SELECT pr.principal_id, pr.name, pr.type_desc,   
    pr.authentication_type_desc, pe.state_desc, pe.permission_name  
FROM sys.database_principals AS pr  
JOIN sys.database_permissions AS pe  
    ON pe.grantee_principal_id = pr.principal_id
where pr.name='{{ .RoleName }}';
{{ end }}

{{ define "addMemberToDBRole" }}
Use {{ .DbName }}
ALTER ROLE {{ .RoleName }} ADD MEMBER {{ .Username }};
{{ end }}


{{ define "dropMemberFromDBRole" }}
Use {{ .DbName }}
ALTER ROLE {{ .RoleName }} DROP MEMBER {{ .Username }};
{{ end }}